package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/pion/udp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake"
	"github.com/willcliffy/kilnwood-game-server/broadcast"
	"github.com/willcliffy/kilnwood-game-server/game"
)

const (
	PORT = 10001
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

	messageBroker := broadcast.NewMessageBroker()
	defer messageBroker.Close()

	gameIdGenerator := sonyflake.NewSonyflake(sonyflake.Settings{})
	gameId, _ := gameIdGenerator.NextID()

	// TODO - support for mulitple games
	newGame := game.NewGame(uint64(gameId), messageBroker)
	newGame.Start()

	// For now, there is one game and it's constantly running
	// TODO - implement lobbies, multiple games, etc
	messageBroker.RegisterGame(gameId, newGame)

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
