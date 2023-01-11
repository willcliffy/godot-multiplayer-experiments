package main

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake"
	"github.com/willcliffy/kilnwood-game-server/broadcast"
	"github.com/willcliffy/kilnwood-game-server/game"
)

const (
	PORT = "9900"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Loggeroni and cheese
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "t"
	log.Logger = log.With().
		Caller().
		Logger()

	messageBroker := broadcast.NewMessageBroker()
	defer messageBroker.Close()

	wsUpgrader := websocket.Upgrader{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Send()
			return
		}
		messageBroker.HandleConnection(c)
	})

	// TODO - support for mulitple games
	gameIdGenerator := sonyflake.NewSonyflake(sonyflake.Settings{})
	gameId, _ := gameIdGenerator.NextID()
	newGame := game.NewGame(uint64(gameId), messageBroker)
	newGame.Start()

	messageBroker.RegisterGame(newGame)

	log.Info().Msgf("Listening on port: %s", PORT)
	err := http.ListenAndServe("localhost:"+PORT, nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("error in ListenAndServe")
	}
}
