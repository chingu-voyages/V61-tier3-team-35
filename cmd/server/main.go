package main

import (
	"log"
	"net/http"
	"os"

	game "github.com/chingu-voyages/V61-tier3-team-35/game"
	"github.com/chingu-voyages/V61-tier3-team-35/internal/api"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	answers, err := game.LoadWords("words/answers.txt")
	if err != nil {
		log.Fatalf("failed to load words: %v", err)
	}
	log.Printf("Loaded %d words", len(answers))

	handler := api.NewHandler(answers)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/health", handlerReadiness)
	mux.HandleFunc("GET /api/daily-word", handler.GetDailyWord)
	mux.HandleFunc("POST /api/guess", handler.EvaluateGuess)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Server started on port: %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
