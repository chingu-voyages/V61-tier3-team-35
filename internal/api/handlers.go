package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chingu-voyages/V61-tier3-team-35/game"
)

type Handler struct {
	answers []string
}

type DailyWord struct {
	Date        string `json:"date"`
	WordLength  int    `json:"word_length"`
	MaxAttempts int    `json:"max_attempts"`
}

type GuessRequest struct {
	Guess string `json:"guess"`
}

func NewHandler(answers []string) *Handler {
	return &Handler{
		answers: answers,
	}
}

func (h *Handler) GetDailyWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"method not allowed", nil)
		return
	}

	dailyWord := game.GetDailyWord(h.answers)
	log.Printf("Today's word: %s", dailyWord)

	today := time.Now().Format(time.DateOnly)
	wordLength := len(dailyWord)
	maxAttempts := 6

	response := DailyWord{
		Date:        today,
		WordLength:  wordLength,
		MaxAttempts: maxAttempts,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) EvaluateGuess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed,
			"method not allowed", nil)
		return
	}

	var guessRequest GuessRequest

	if err := json.NewDecoder(r.Body).Decode(&guessRequest); err != nil {
		respondWithError(w, http.StatusBadRequest,
			"Failed to decode request body", err)
		return
	}

	if guessRequest.Guess == "" || len(guessRequest.Guess) > 5 {
		respondWithError(w, http.StatusBadRequest, "length of guess should be '5'", nil)
		return
	}

	dailyWord := game.GetDailyWord(h.answers)
	result := game.EvaluateGuess(guessRequest.Guess, dailyWord)

	respondWithJSON(w, http.StatusOK, result)

}
