package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/willcliffy/kilnwood-game-server/game/objects"
)

const bufSize = 8192

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "kilnwood-game.com", Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		//resBytes, err := io.ReadAll(res.Body)
		//log.Printf("res: %v %v", err, string(resBytes))
		log.Fatal("dial: ", err)
	}

	defer conn.Close()

	fmt.Println("Connected; type 'exit' to shutdown gracefully")

	id := ""
	// Simulate a chat session
	go func() {
		b := make([]byte, bufSize)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				panic(err)
			}
			print(string(msg))
			if len(id) == 0 {
				if string(b[0]) == "d" {
					fmt.Println("Server disconnected")
				}

				x := struct {
					Type     string
					PlayerId string
					Spawn    objects.Location
				}{}

				err := json.Unmarshal(msg, &x)
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

		if err = conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			panic(err)
		}
	}
}
