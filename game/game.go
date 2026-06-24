package game

import "strings"

type LetterFeedback struct {
	Letter string `json:"letter"`
	Status string `json:"status"`
}

type GuessResult struct {
	Feedback  []LetterFeedback `json:"feedback"`
	IsCorrect bool             `json:"is_correct"`
}

const (
	Correct = "correct"
	Present = "present"
	Absent  = "absent"
)

func EvaluateGuess(guess, target string) GuessResult {
	guess = strings.ToLower(guess)
	target = strings.ToLower(target)

	feedback := make([]LetterFeedback, len(guess))
	remaining := make(map[rune]int)

	guessRunes := []rune(guess)
	targetRunes := []rune(target)

	isCorrect := guess == target

	for i := range guessRunes {
		feedback[i] = LetterFeedback{
			Letter: string(guessRunes[i]),
			Status: Absent,
		}

		if guessRunes[i] == targetRunes[i] {
			feedback[i].Status = Correct
		} else {
			remaining[targetRunes[i]]++
		}
	}

	for i := range guessRunes {
		if feedback[i].Status == Correct {
			continue
		}

		if remaining[guessRunes[i]] > 0 {
			feedback[i].Status = Present
			remaining[guessRunes[i]]--
		}
	}

	return GuessResult{
		Feedback:  feedback,
		IsCorrect: isCorrect,
	}
}
