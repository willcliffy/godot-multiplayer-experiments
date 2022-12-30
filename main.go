package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/game"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Loggeroni and cheese
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	zerolog.CallerFieldName = "c"
	log.Logger = log.With().Caller().Logger()

	// Prepare the configuration of the DTLS connection
	cert, certPool := loadCertAndCertPool()
	config := &dtls.Config{
		Certificates:         []tls.Certificate{cert},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		ClientAuth:           dtls.RequireAndVerifyClientCert,
		ClientCAs:            certPool,
		// Create timeout context for accepted connection.
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(ctx, 30*time.Second)
		},
	}

	// Connect to a DTLS server
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 4444}
	listener, err := dtls.Listen("udp", addr, config)
	if err != nil {
		log.Fatal().Err(err)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatal().Err(err)
		}
	}()

	log.Info().Msg("Listening")

	// For now, there is one game and it's constantly running
	// TODO - implement lobbies, multiple games, etc
	game := game.NewGame()
	game.Start()
	defer game.Stop()

	messageBroker := NewMessageBroker()
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal().Err(err)
			} else if conn == nil {
				continue
			}

			dtlsConn, ok := conn.(*dtls.Conn)
			if !ok {
				log.Fatal().Msgf("Connection to %v was not DTLS!", conn)
			}

			// It is the messageBroker's responsibility to close the connection
			messageBroker.RegisterClient(dtlsConn)
		}
	}()

	// Gracefully handle shutdown signal, block until done
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
}

func loadCertAndCertPool() (tls.Certificate, *x509.CertPool) {
	certificate, err := tls.LoadX509KeyPair(
		"cert/server.pub.pem",
		"cert/server.pem")
	if err != nil {
		log.Fatal().Err(err)
	}

	rawRootCertData, err := os.ReadFile(filepath.Clean("cert/server.pub.pem"))
	if err != nil {
		log.Fatal().Err(err)
	}

	var rootCertificate tls.Certificate
	for {
		block, rest := pem.Decode(rawRootCertData)
		if block == nil {
			break
		}

		if block.Type != "CERTIFICATE" {
			log.Fatal().Msg("block is not cert")
		}

		rootCertificate.Certificate = append(rootCertificate.Certificate, block.Bytes)
		rawRootCertData = rest
	}

	if len(rootCertificate.Certificate) == 0 {
		log.Fatal().Msg("no cert found")
	}

	certPool := x509.NewCertPool()
	cert, err := x509.ParseCertificate(rootCertificate.Certificate[0])
	if err != nil {
		log.Fatal().Err(err)
	}

	certPool.AddCert(cert)

	return certificate, certPool
}
