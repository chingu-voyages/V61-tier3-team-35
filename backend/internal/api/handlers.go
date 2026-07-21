package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chingu-voyages/V61-tier3-team-35/backend/game"
)

const (
	StatusInProgress = "in_progress"
	StatusWon        = "won"
	StatusLost       = "lost"

	MaxAttempts    = 6
	ClientIDCookie = "client_id"
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

type GameState struct {
	AttemptsUsed int           `json:"attempts_used"`
	Status       string        `json:"status"`
	Guesses      []GuessResult `json:"guesses"`
}

type GameResponse struct {
	Date         string        `json:"date,omitempty"`
	WordLength   int           `json:"word_length"`
	MaxAttempts  int           `json:"max_attempts"`
	Status       string        `json:"status"`
	AttemptsUsed int           `json:"attempts_used"`
	Guesses      []GuessResult `json:"guesses"`
}

type GuessRequest struct {
	Guess string `json:"guess"`
}

type GuessResult struct {
	Word     string `json:"word"`
	Feedback any    `json:"feedback"`
}

type GuessResponse struct {
	Feedback     any           `json:"feedback"`
	IsCorrect    bool          `json:"is_correct"`
	Status       string        `json:"status"`
	AttemptsUsed int           `json:"attempts_used"`
	TargetWord   string        `json:"target_word,omitempty"`
	Guesses      []GuessResult `json:"guesses"`
}

func NewHandler(answers []string, validWords []string, production bool) *Handler {
	validGuesses := make(map[string]struct{}, len(validWords))

	for _, word := range validWords {
		validGuesses[strings.ToLower(word)] = struct{}{}
	}

	return &Handler{
		answers:      answers,
		validGuesses: validGuesses,
		games:        make(map[string]PlayerGames),
		production:   production,
	}
}

func newGameState() GameState {
	return GameState{
		AttemptsUsed: 0,
		Status:       StatusInProgress,
		Guesses:      []GuessResult{},
	}
}

func newDailyGame(date string) DailyGame {
	return DailyGame{
		Date:  date,
		State: newGameState(),
	}
}

func newPracticeGame(targetWord string) PracticeGame {
	return PracticeGame{
		TargetWord: targetWord,
		State:      newGameState(),
	}
}

func newGameResponse(date string, wordLength int, state GameState) GameResponse {
	return GameResponse{
		Date:         date,
		WordLength:   wordLength,
		MaxAttempts:  MaxAttempts,
		Status:       state.Status,
		AttemptsUsed: state.AttemptsUsed,
		Guesses:      state.Guesses,
	}
}

func generateClientID() string {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return time.Now().Format("20060102150405.000000000")
	}

	return hex.EncodeToString(bytes)
}

func (h *Handler) getOrSetClientID(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie(ClientIDCookie)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	clientID := generateClientID()

	secure := false
	sameSite := http.SameSiteLaxMode

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

func (h *Handler) GetOrCreateDailyGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	clientID := h.getOrSetClientID(w, r)
	today := time.Now().Format(time.DateOnly)

	h.mu.Lock()
	dailyGame := h.getOrCreateDailyGame(clientID, today)
	h.mu.Unlock()

	dailyWord := game.GetDailyWord(h.answers)

	response := newGameResponse(
		today,
		len(dailyWord),
		dailyGame.State,
	)

	respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) getOrCreateDailyGame(clientID string, today string) DailyGame {
	playerGames, ok := h.games[clientID]

	if !ok || playerGames.DailyGame.Date != today {
		playerGames.DailyGame = newDailyGame(today)
		h.games[clientID] = playerGames
	}

	return playerGames.DailyGame
}

func (h *Handler) savePlayerDailyGame(clientID string, dailyGame DailyGame) {
	playerGames := h.games[clientID]
	playerGames.DailyGame = dailyGame
	h.games[clientID] = playerGames
}

func (h *Handler) getPracticeGame(clientID string) (PracticeGame, bool) {
	playerGames, ok := h.games[clientID]

	if !ok || playerGames.PracticeGame.TargetWord == "" {
		return PracticeGame{}, false
	}

	return playerGames.PracticeGame, true
}

func (h *Handler) getOrCreatePracticeGame(clientID string) PracticeGame {
	if practiceGame, ok := h.getPracticeGame(clientID); ok {
		return practiceGame
	}

	targetWord := game.GetRandomWord(h.answers)

	return h.startPracticeGame(clientID, targetWord)
}

// startPracticeGame intentionally replaces any existing practice game.
func (h *Handler) startPracticeGame(clientID string, targetWord string) PracticeGame {
	playerGames := h.games[clientID]
	playerGames.PracticeGame = newPracticeGame(targetWord)
	h.games[clientID] = playerGames

	return playerGames.PracticeGame
}

func (h *Handler) savePlayerPracticeGame(clientID string, practiceGame PracticeGame) {
	playerGames := h.games[clientID]
	playerGames.PracticeGame = practiceGame
	h.games[clientID] = playerGames
}

func (h *Handler) hasCompletedDailyGame(clientID string, today string) bool {
	playerGames, ok := h.games[clientID]
	if !ok {
		return false
	}

	dailyGame := playerGames.DailyGame

	if dailyGame.Date != today {
		return false
	}

	return dailyGame.State.Status == StatusWon ||
		dailyGame.State.Status == StatusLost
}

// GetOrCreatePracticeGame returns the current practice game.
// It creates the player's first practice game only when none exists.
func (h *Handler) GetOrCreatePracticeGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	clientID := h.getOrSetClientID(w, r)
	today := time.Now().Format(time.DateOnly)

	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.hasCompletedDailyGame(clientID, today) {
		respondWithError(
			w,
			http.StatusForbidden,
			"complete the daily game before accessing practice mode",
			nil,
		)
		return
	}

	practiceGame := h.getOrCreatePracticeGame(clientID)

	response := newGameResponse(
		"",
		len(practiceGame.TargetWord),
		practiceGame.State,
	)

	respondWithJSON(w, http.StatusOK, response)
}

// StartNewPracticeGame intentionally discards the current practice game
// and creates a new one.
func (h *Handler) StartNewPracticeGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	clientID := h.getOrSetClientID(w, r)
	today := time.Now().Format(time.DateOnly)

	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.hasCompletedDailyGame(clientID, today) {
		respondWithError(
			w,
			http.StatusForbidden,
			"complete the daily game before accessing practice mode",
			nil,
		)
		return
	}

	targetWord := game.GetRandomWord(h.answers)
	practiceGame := h.startPracticeGame(clientID, targetWord)

	response := newGameResponse(
		"",
		len(practiceGame.TargetWord),
		practiceGame.State,
	)

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *Handler) EvaluateGuess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request GuessRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"failed to decode request body",
			err,
		)
		return
	}

	guess := strings.ToLower(strings.TrimSpace(request.Guess))

	if len(guess) != 5 {
		respondWithError(
			w,
			http.StatusBadRequest,
			"guess must be exactly 5 letters",
			nil,
		)
		return
	}

	if _, ok := h.validGuesses[guess]; !ok {
		respondWithError(
			w,
			http.StatusBadRequest,
			"word is not in the accepted word list",
			nil,
		)
		return
	}

	clientID := h.getOrSetClientID(w, r)
	today := time.Now().Format(time.DateOnly)

	h.mu.Lock()
	defer h.mu.Unlock()

	switch r.URL.Path {
	case "/api/daily/guess":
		h.evaluateDailyGuess(w, clientID, today, guess)

	case "/api/practice/guess":
		h.evaluatePracticeGuess(w, clientID, today, guess)

	default:
		respondWithError(
			w,
			http.StatusNotFound,
			"game route not found",
			nil,
		)
	}
}

func (h *Handler) evaluateDailyGuess(w http.ResponseWriter, clientID string, today string, guess string) {
	dailyGame := h.getOrCreateDailyGame(clientID, today)

	if dailyGame.State.Status != StatusInProgress {
		respondWithError(
			w,
			http.StatusConflict,
			"game is already finished",
			nil,
		)
		return
	}

	targetWord := strings.ToLower(game.GetDailyWord(h.answers))
	result := game.CompareGuessToTarget(guess, targetWord)

	dailyGame.State.AttemptsUsed++
	dailyGame.State.Guesses = append(
		dailyGame.State.Guesses,
		GuessResult{
			Word:     guess,
			Feedback: result.Feedback,
		},
	)

	response := GuessResponse{
		Feedback:     result.Feedback,
		IsCorrect:    result.IsCorrect,
		Status:       StatusInProgress,
		AttemptsUsed: dailyGame.State.AttemptsUsed,
		Guesses:      dailyGame.State.Guesses,
	}

	if result.IsCorrect {
		dailyGame.State.Status = StatusWon
		response.Status = StatusWon
	} else if dailyGame.State.AttemptsUsed >= MaxAttempts {
		dailyGame.State.Status = StatusLost
		response.Status = StatusLost
		response.TargetWord = targetWord
	}

	h.savePlayerDailyGame(clientID, dailyGame)

	respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) evaluatePracticeGuess(w http.ResponseWriter, clientID string, today string, guess string) {
	if !h.hasCompletedDailyGame(clientID, today) {
		respondWithError(
			w,
			http.StatusForbidden,
			"complete the daily game before accessing practice mode",
			nil,
		)
		return
	}

	practiceGame, ok := h.getPracticeGame(clientID)
	if !ok {
		respondWithError(
			w,
			http.StatusConflict,
			"practice game has not been started",
			nil,
		)
		return
	}

	if practiceGame.State.Status != StatusInProgress {
		respondWithError(
			w,
			http.StatusConflict,
			"practice game is already finished",
			nil,
		)
		return
	}

	targetWord := strings.ToLower(practiceGame.TargetWord)
	result := game.CompareGuessToTarget(guess, targetWord)

	practiceGame.State.AttemptsUsed++
	practiceGame.State.Guesses = append(
		practiceGame.State.Guesses,
		GuessResult{
			Word:     guess,
			Feedback: result.Feedback,
		},
	)

	response := GuessResponse{
		Feedback:     result.Feedback,
		IsCorrect:    result.IsCorrect,
		Status:       StatusInProgress,
		AttemptsUsed: practiceGame.State.AttemptsUsed,
		Guesses:      practiceGame.State.Guesses,
	}

	if result.IsCorrect {
		practiceGame.State.Status = StatusWon
		response.Status = StatusWon
	} else if practiceGame.State.AttemptsUsed >= MaxAttempts {
		practiceGame.State.Status = StatusLost
		response.Status = StatusLost
		response.TargetWord = targetWord
	}

	h.savePlayerPracticeGame(clientID, practiceGame)

	respondWithJSON(w, http.StatusOK, response)
}
