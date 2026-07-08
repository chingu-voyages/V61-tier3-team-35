package main

import (
	"log"
	"net/http"
	"os"

	game "github.com/chingu-voyages/V61-tier3-team-35/game"
	"github.com/chingu-voyages/V61-tier3-team-35/internal/api"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	Production bool
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	env := os.Getenv("ENV") == "production"

	if port == "" {
		port = "8080"
	}

	cfg := Config{
		Port:       port,
		Production: env,
	}

	validWords, err := game.LoadWords("words/allowed-guess.txt")
	if err != nil {
		log.Fatalf("failed to load valid words: %v", err)
	}

	answers, err := game.LoadWords("words/answers.txt")
	if err != nil {
		log.Fatalf("failed to load answers: %v", err)
	}

	handler := api.NewHandler(answers, validWords, cfg.Production)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/health", handlerReadiness)
	mux.HandleFunc("GET /api/daily-word", handler.GetDailyWord)
	mux.HandleFunc("POST /api/guess", handler.EvaluateGuess)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: api.CorsMiddleware(mux),
	}
	log.Printf("Server started on port: %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
