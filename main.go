package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

/*type apiConfig struct {
	baseURL string
}*/

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	/*	_ = &apiConfig{
		baseURL: os.Getenv("BASE_URL"),
	}*/

	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/health", handlerReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Server started on port: %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
