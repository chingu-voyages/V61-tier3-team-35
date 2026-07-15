package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chingu-voyages/V61-tier3-team-35/game"
)

const (
	StatusInProgress = "in_progress"
	StatusWon        = "won"
	StatusLost       = "lost"
	MaxAttempts      = 6
	ClientIDCookie   = "client_id"
)

type Handler struct {
	answers      []string
	validGuesses map[string]struct{}
	mu           sync.Mutex
	games        map[string]PlayerGames
	production   bool
}

type PlayerGames struct {
	DailyGame    DailyGame
	PracticeGame PracticeGame
}

type DailyGame struct {
	Date  string
	State GameState
}

type PracticeGame struct {
	TargetWord string
	State      GameState
}

type DailyWord struct {
	Date         string   `json:"date"`
	WordLength   int      `json:"word_length"`
	MaxAttempts  int      `json:"max_attempts"`
	Status       string   `json:"status"`
	AttemptsUsed int      `json:"attempts_used"`
	Guesses      []string `json:"guesses"`
}

type GuessRequest struct {
	Guess string `json:"guess"`
}

type GameState struct {
	AttemptsUsed int      `json:"attempts_used"`
	Status       string   `json:"status"`
	Guesses      []string `json:"guesses"`
}

type GuessResponse struct {
	Feedback     any      `json:"feedback"`
	IsCorrect    bool     `json:"is_correct"`
	Status       string   `json:"status"`
	AttemptsUsed int      `json:"attempts_used"`
	TargetWord   string   `json:"target_word,omitempty"`
	Guesses      []string `json:"guesses,omitempty"`
}

func NewHandler(answers []string, validWords []string, env bool) *Handler {
	validGuesses := make(map[string]struct{})

	for _, word := range validWords {
		validGuesses[strings.ToLower(word)] = struct{}{}
	}

	return &Handler{
		answers:      answers,
		validGuesses: validGuesses,
		games:        make(map[string]PlayerGames),
		production:   env,
	}
}

func newPracticeGame(targetWord string) PracticeGame {
	return PracticeGame{
		TargetWord: targetWord,
		State: GameState{
			AttemptsUsed: 0,
			Status:       StatusInProgress,
			Guesses:      []string{},
		},
	}
}

func newDailyGame(date string) DailyGame {
	return DailyGame{
		Date: date,
		State: GameState{
			AttemptsUsed: 0,
			Status:       StatusInProgress,
			Guesses:      []string{},
		},
	}
}

func (h *Handler) getOrCreateDailyGame(clientID, today string) DailyGame {
	playerGame, ok := h.games[clientID]
	if !ok || playerGame.DailyGame.Date != today {
		playerGame.DailyGame = newDailyGame(today)
		h.games[clientID] = playerGame
	}

	return playerGame.DailyGame
}

func (h *Handler) savePlayerDailyGame(clientID string, dailyGame DailyGame) {
	playerGame, ok := h.games[clientID]
	if ok {
		playerGame.DailyGame = dailyGame
		h.games[clientID] = playerGame
	} else {
		h.games[clientID] = PlayerGames{
			DailyGame: dailyGame,
		}
	}
}

func (h *Handler) getOrSetClientID(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie(ClientIDCookie)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	clientID := generateClientID()
	sameSite := http.SameSiteLaxMode
	secure := false
	if h.production {
		secure = true
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     ClientIDCookie,
		Value:    clientID,
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
		MaxAge:   60 * 60 * 24 * 30,
	})

	return clientID
}

func generateClientID() string {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return time.Now().Format("20060102150405.000000000")
	}

	return hex.EncodeToString(bytes)
}

func (h *Handler) GetDailyWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	clientID := h.getOrSetClientID(w, r)
	today := time.Now().Format(time.DateOnly)

	h.mu.Lock()
	playerGame := h.getOrCreateDailyGame(clientID, today)
	h.mu.Unlock()

	dailyWord := game.GetDailyWord(h.answers)
	log.Printf("Today's word: %s", dailyWord)

	response := DailyWord{
		Date:         today,
		WordLength:   len(dailyWord),
		MaxAttempts:  MaxAttempts,
		Status:       playerGame.State.Status,
		AttemptsUsed: playerGame.State.AttemptsUsed,
		Guesses:      playerGame.State.Guesses,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) EvaluateGuess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed,
			"method not allowed", nil)
		return
	}

	var guessRequest GuessRequest

	if err := json.NewDecoder(r.Body).Decode(&guessRequest); err != nil {
		respondWithError(w, http.StatusBadRequest,
			"failed to decode request body", err)
		return
	}

	guess := strings.ToLower(strings.TrimSpace(guessRequest.Guess))

	if len(guess) != 5 {
		respondWithError(w, http.StatusBadRequest,
			"guess must be exactly 5 letters", nil)
		return
	}

	if _, ok := h.validGuesses[guess]; !ok {
		respondWithError(w, http.StatusBadRequest,
			"word is not in the accepted word list", nil)
		return
	}

	clientID := h.getOrSetClientID(w, r)
	today := time.Now().Format(time.DateOnly)

	h.mu.Lock()
	defer h.mu.Unlock()

	playerGame := h.getOrCreateDailyGame(clientID, today)

	if playerGame.State.Status != StatusInProgress {
		respondWithError(w, http.StatusConflict,
			"game is already finished", nil)
		return
	}

	dailyWord := strings.ToLower(game.GetDailyWord(h.answers))
	result := game.EvaluateGuess(guess, dailyWord)

	playerGame.State.AttemptsUsed++
	playerGame.State.Guesses = append(playerGame.State.Guesses, guess)

	response := GuessResponse{
		Feedback:     result.Feedback,
		IsCorrect:    result.IsCorrect,
		Status:       StatusInProgress,
		AttemptsUsed: playerGame.State.AttemptsUsed,
		Guesses:      playerGame.State.Guesses,
	}

	if guess == dailyWord {
		playerGame.State.Status = StatusWon
		response.Status = StatusWon
	} else if playerGame.State.AttemptsUsed >= MaxAttempts {
		playerGame.State.Status = StatusLost
		response.Status = StatusLost
		response.TargetWord = dailyWord
	}

	h.savePlayerDailyGame(clientID, playerGame)

	respondWithJSON(w, http.StatusOK, response)
}
