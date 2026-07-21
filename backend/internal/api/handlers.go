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

	"github.com/chingu-voyages/V61-tier3-team-35/backend/game"
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

type GameState struct {
	AttemptsUsed int           `json:"attempts_used"`
	Status       string        `json:"status"`
	Guesses      []GuessResult `json:"guesses"`
}

type GuessResponse struct {
	Feedback     any           `json:"feedback"`
	IsCorrect    bool          `json:"is_correct"`
	Status       string        `json:"status"`
	AttemptsUsed int           `json:"attempts_used"`
	TargetWord   string        `json:"target_word"`
	Guesses      []GuessResult `json:"guesses"`
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
			Guesses:      []GuessResult{},
		},
	}
}

func newDailyGame(date string) DailyGame {
	return DailyGame{
		Date: date,
		State: GameState{
			AttemptsUsed: 0,
			Status:       StatusInProgress,
			Guesses:      []GuessResult{},
		},
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

func (h *Handler) GetOrCreateDailyGame(w http.ResponseWriter, r *http.Request) {
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

	response := GameResponse{
		Date:         today,
		WordLength:   len(dailyWord),
		MaxAttempts:  MaxAttempts,
		Status:       playerGame.State.Status,
		AttemptsUsed: playerGame.State.AttemptsUsed,
		Guesses:      playerGame.State.Guesses,
	}

	respondWithJSON(w, http.StatusOK, response)
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

func (h *Handler) createPracticeGame(clientID, targetWord string) PracticeGame {
	playerGames := h.games[clientID]

	playerGames.PracticeGame = newPracticeGame(targetWord)
	h.games[clientID] = playerGames

	return playerGames.PracticeGame
}

func (h *Handler) getPracticeGame(clientID string) (PracticeGame, bool) {
	playerGames, ok := h.games[clientID]
	if !ok || playerGames.PracticeGame.TargetWord == "" {
		return PracticeGame{}, false
	}

	return playerGames.PracticeGame, true
}

func (h *Handler) savePlayerPracticeGame(
	clientID string,
	practiceGame PracticeGame,
) {
	playerGames := h.games[clientID]
	playerGames.PracticeGame = practiceGame
	h.games[clientID] = playerGames
}

func (h *Handler) GetOrCreatePracticeGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	clientID := h.getOrSetClientID(w, r)
	today := time.Now().Format(time.DateOnly)

	h.mu.Lock()
	playerGame := h.getOrCreateDailyGame(clientID, today)
	defer h.mu.Unlock()

	if playerGame.State.Status == StatusInProgress {
		respondWithError(w, http.StatusForbidden,
			"cannot start a new practice game while daily game is in progress", nil)
		return
	}

	randomWord := game.GetRandomWord(h.answers)
	newPracticeGame := h.createPracticeGame(clientID, randomWord)

	response := GameResponse{
		WordLength:   len(randomWord),
		MaxAttempts:  MaxAttempts,
		Status:       newPracticeGame.State.Status,
		AttemptsUsed: newPracticeGame.State.AttemptsUsed,
		Guesses:      newPracticeGame.State.Guesses,
	}

	respondWithJSON(w, http.StatusOK, response)

}

func (h *Handler) EvaluateGuess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gameMode := r.URL.Path
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

	switch gameMode {
	case "/api/guess":
		dailyGame := h.getOrCreateDailyGame(clientID, today)

		if dailyGame.State.Status != StatusInProgress {
			respondWithError(w, http.StatusConflict,
				"game is already finished", nil)
			return
		}

		dailyWord := strings.ToLower(game.GetDailyWord(h.answers))
		result := game.CompareGuessToTarget(guess, dailyWord)

		dailyGame.State.AttemptsUsed++

		guessResult := GuessResult{
			Word:     guess,
			Feedback: result.Feedback,
		}

		dailyGame.State.Guesses = append(
			dailyGame.State.Guesses,
			guessResult,
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
			response.TargetWord = dailyWord
		}

		h.savePlayerDailyGame(clientID, dailyGame)

		respondWithJSON(w, http.StatusOK, response)

	case "/api/practice/guess":
		practiceGame, ok := h.getPracticeGame(clientID)
		if !ok {
			respondWithError(w, http.StatusConflict,
				"practice game has not been started",
				nil,
			)
			return
		}

		if practiceGame.State.Status != StatusInProgress {
			respondWithError(w, http.StatusConflict,
				"practice game is already finished",
				nil,
			)
			return
		}

		targetWord := strings.ToLower(practiceGame.TargetWord)
		result := game.CompareGuessToTarget(guess, targetWord)

		practiceGame.State.AttemptsUsed++

		guessResult := GuessResult{
			Word:     guess,
			Feedback: result.Feedback,
		}

		practiceGame.State.Guesses = append(
			practiceGame.State.Guesses,
			guessResult,
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

	default:
		respondWithError(w, http.StatusNotFound,
			"game route not found",
			nil,
		)
	}
}
