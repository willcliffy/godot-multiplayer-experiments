package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/willcliffy/kilnwood-game-server/game/objects"
)

const bufSize = 8192

func main() {
	// Prepare the IP to connect to
	addr := &net.UDPAddr{IP: net.ParseIP("35.227.75.95"), Port: 10001}

	udpConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := udpConn.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Connected; type 'exit' to shutdown gracefully")

	id := ""
	// Simulate a chat session
	go func() {
		b := make([]byte, bufSize)

		for {
			n, err := udpConn.Read(b)
			if err != nil {
				panic(err)
			}
			if len(id) == 0 {
				if string(b[0]) == "d" {
					fmt.Println("Server disconnected")
				}

				x := struct {
					Type     string
					PlayerId string
					Spawn    objects.Position
				}{}

				err := json.Unmarshal(b[:n], &x)
				if err != nil {
					panic(err)
				}
				id = x.PlayerId
				fmt.Printf("id: %v", id)
			}
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

		_, err = udpConn.Write([]byte(text))
		if err != nil {
			panic(err)
		}
	}
}
