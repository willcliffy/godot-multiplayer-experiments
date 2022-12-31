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

const (
	PORT = 4444
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
	log.Logger = log.With().
		Caller().
		Logger()

	// Set up Cert and Cert Pool - this is used in our DTLS connections
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

	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: PORT}
	listener, err := dtls.Listen("udp", addr, config)
	if err != nil {
		log.Fatal().Err(err)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatal().Err(err)
		}
	}()

	log.Info().Msgf("Listening on port: %d", PORT)

	// TODO - how does message broker communicate with game(s)

	messageBroker := NewMessageBroker()
	defer messageBroker.Close()

	// For now, there is one game and it's constantly running
	// TODO - implement lobbies, multiple games, etc
	game := game.NewGame()
	game.Start()
	defer game.Stop()

	messageBroker.RegisterGame(game)

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

			// We pass off the connection to messagebroker, who is now responsible for closing it
			messageBroker.RegisterConnection(dtlsConn)
		}
	}()

	// Gracefully handle shutdown signal, block until done
	shutdown := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		close(shutdown)
	}()
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
