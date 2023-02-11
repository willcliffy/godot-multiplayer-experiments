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
	PORT              string = "8080"
	WEBCLIENT_ENABLED bool   = true
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

	// TODO - support for mulitple games
	gameIdGenerator := sonyflake.NewSonyflake(sonyflake.Settings{})
	gameId, _ := gameIdGenerator.NextID()
	newGame := game.NewGame(gameId, messageBroker)
	newGame.Start()
	defer newGame.Close()

	messageBroker.RegisterMessageReceiver(newGame)

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: ConfigureRouter(messageBroker),
	}

	log.Info().Msgf("Listening on port: %s", PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msgf("error in ListenAndServe")
	}
}

func ConfigureRouter(mb *broadcast.MessageBroker) http.Handler {
	router := chi.NewRouter()

	// Health checks
	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("alive\n"))
	})

	// Web client
	if WEBCLIENT_ENABLED {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dist/index.html")
		})
		router.Get("/{file}", func(w http.ResponseWriter, r *http.Request) {
			f := chi.URLParam(r, "file")
			http.ServeFile(w, r, "dist/"+f)
		})
	}

	// Game server - websockets
	wsUpgrader := websocket.Upgrader{}
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true } // TODO - security risk
	router.Route("/ws", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/connect", func(w http.ResponseWriter, r *http.Request) {
				c, err := wsUpgrader.Upgrade(w, r, nil)
				if err != nil {
					log.Error().Err(err).Send()
					return
				}

				mb.RegisterAndHandleWebsocketConnection(c)
			})
		})
	})

	return router
}
