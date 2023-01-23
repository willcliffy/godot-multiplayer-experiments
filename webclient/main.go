package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PORT = "3000"
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

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("alive\n"))
	})

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("alive\n"))
	})

	router.Get("/game/", GetServeStaticFile("index.html"))
	router.Get("/game/{filename}", GetServeStaticFiles)

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}

	log.Info().Msgf("Listening on port: %s", PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msgf("error in ListenAndServe")
	}
}

func GetServeStaticFile(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("dist/%v", filename))
	}
}

func GetServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	http.ServeFile(w, r, fmt.Sprintf("dist/%v", filename))
}
