package api

import (
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

func NewHandler(answers []string) *Handler {
	return &Handler{
		answers: answers,
	}
}

func (h *Handler) HandlerDailyWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed,
			"method not allowed", nil)
		return
	}

	dailyWord := game.GetDailyWord(h.answers)
	today := time.Now().String()
	wordLength := len(dailyWord)
	maxAttempts := 6

	response := DailyWord{
		Date:        today,
		WordLength:  wordLength,
		MaxAttempts: maxAttempts,
	}

	respondWithJSON(w, http.StatusOK, response)
}
