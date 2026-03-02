package main

import (
	"net/http"

	"github.com/fevziatanoglu/test-go-project/internal/config"
	"github.com/fevziatanoglu/test-go-project/internal/middleware"
	"github.com/fevziatanoglu/test-go-project/internal/router"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.LoadConfig()

	r := router.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.Logging,
		middleware.CORS(cfg),
		middleware.SecurityHeaders,
		middleware.RateLimit(cfg.RateLimitRPS, cfg.RateLimitBurst),
	)

	r.Handle("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	addr := ":" + cfg.ServerPort
	log.Info().Str("addr", addr).Msg("http_server_starting")
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal().Err(err).Msg("http_server_error")
	}
}
