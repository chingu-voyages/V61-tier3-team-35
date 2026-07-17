package game

import (
	"testing"
)

func TestAnswersAreSubsetOfAllowedGuesses(t *testing.T) {
	allowed, err := LoadWords("../words/allowed-guess.txt")
	if err != nil {
		t.Fatalf("failed to load allowed guesses: %v", err)
	}

	answers, err := LoadWords("../words/answers.txt")
	if err != nil {
		t.Fatalf("failed to load answers: %v", err)
	}

	allowedSet := make(map[string]struct{}, len(allowed))
	for _, word := range allowed {
		allowedSet[word] = struct{}{}
	}

	for _, answer := range answers {
		if _, ok := allowedSet[answer]; !ok {
			t.Errorf("answer %q is missing from allowed-guess.txt", answer)
		}
	}
}
