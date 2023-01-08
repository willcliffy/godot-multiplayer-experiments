package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/pion/udp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/game"
)

const (
	PORT = 4444
)

func main() {
	_, cancel := context.WithCancel(context.Background())
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

	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: PORT}
	listener, err := udp.Listen("udp", addr)
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
	game := game.NewGame(messageBroker)
	game.Start()
	defer game.Stop()

	messageBroker.RegisterGame(game)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Error().Err(err)
				continue
			} else if conn == nil {
				continue
			}

			// We pass off the connection to messagebroker, who is now responsible for closing it
			messageBroker.RegisterConnection(conn)
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
