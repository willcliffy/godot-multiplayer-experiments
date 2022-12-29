package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/willcliffy/kilnwood-game-server/hub"
)

func main() {
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 4444}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Everything below is the pion-DTLS API
	certificate, err := tls.LoadX509KeyPair(
		"cert/server.pub.pem",
		"cert/server.pem")
	if err != nil {
		panic(err)
	}

	rootCertificate, err := LoadCertificate("cert/server.pub.pem")
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	cert, err := x509.ParseCertificate(rootCertificate.Certificate[0])
	if err != nil {
		panic(err)
	}

	certPool.AddCert(cert)

	// Prepare the configuration of the DTLS connection
	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		ClientAuth:           dtls.RequireAndVerifyClientCert,
		ClientCAs:            certPool,
		// Create timeout context for accepted connection.
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(ctx, 30*time.Second)
		},
	}

	// Connect to a DTLS server
	listener, err := dtls.Listen("udp", addr, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Listening")

	// Simulate a chat session
	hub := hub.NewHub()

	go func() {
		for {
			// Wait for a connection.
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}

			// defer conn.Close() // TODO: graceful shutdown

			// conn is of type net.Conn but can be cast to dtls.Conn:
			// dtlsConn := conn.(*dtls.Conn)

			// Register the connection with the chat hub
			hub.Register(conn)
		}
	}()

	// Start chatting
	hub.Chat()
}

func LoadCertificate(path string) (*tls.Certificate, error) {
	rawData, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	var certificate tls.Certificate

	for {
		block, rest := pem.Decode(rawData)
		if block == nil {
			break
		}

		if block.Type != "CERTIFICATE" {
			panic("block is not cert")
		}

		certificate.Certificate = append(certificate.Certificate, block.Bytes)
		rawData = rest
	}

	if len(certificate.Certificate) == 0 {
		panic("no cert found")
	}

	return &certificate, nil
}
