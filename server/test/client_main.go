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
	addr, err := net.ResolveUDPAddr("udp", "localhost:9900")
	if err != nil {
		panic(err)
	}

	udpConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}

	defer udpConn.Close()

	fmt.Println("Connected; type 'exit' to shutdown gracefully")

	id := ""
	// Simulate a chat session
	go func() {
		b := make([]byte, bufSize)

		for {
			n, _, err := udpConn.ReadFromUDP(b)
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
