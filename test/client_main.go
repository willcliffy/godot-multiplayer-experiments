package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/examples/util"
)

const bufSize = 8192

func main() {
	// Prepare the IP to connect to
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 4444}

	// Everything below is the pion-DTLS API!
	certificate, err := tls.LoadX509KeyPair(
		"cert/client.pub.pem",
		"cert/client.pem")
	if err != nil {
		panic(err)
	}

	rootCertificate, err := util.LoadCertificate("cert/server.pub.pem")
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
		RootCAs:              certPool,
	}

	// Connect to a DTLS server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dtlsConn, err := dtls.DialWithContext(ctx, "udp", addr, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := dtlsConn.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Connected; type 'exit' to shutdown gracefully")

	// Simulate a chat session
	go func() {
		b := make([]byte, bufSize)

		for {
			n, err := dtlsConn.Read(b)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Got message: %s\n", string(b[:n]))
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if strings.TrimSpace(text) == "exit" {
			return
		}

		_, err = dtlsConn.Write([]byte(text))
		if err != nil {
			panic(err)
		}
	}
}
