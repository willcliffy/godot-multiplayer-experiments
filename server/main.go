package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake"
	"github.com/willcliffy/kilnwood-game-server/broadcast"
	"github.com/willcliffy/kilnwood-game-server/game"
)

const (
	PORT = "8080"
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

	// TODO - This is a security risk and should be patched before prod
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("alive\n"))
	})

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("alive\n"))
	})

	router.Get("/connect", func(w http.ResponseWriter, r *http.Request) {
		c, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Send()
			return
		}
		messageBroker.RegisterAndHandleWebsocketConnection(c)
	})

	// TODO - support for mulitple games
	gameIdGenerator := sonyflake.NewSonyflake(sonyflake.Settings{})
	gameId, _ := gameIdGenerator.NextID()
	newGame := game.NewGame(broadcast.GameId(gameId), messageBroker)
	newGame.Start()

	messageBroker.RegisterMessageReceiver(newGame)

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}

	log.Info().Msgf("Listening on port: %s", PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msgf("error in ListenAndServe")
	}
}
