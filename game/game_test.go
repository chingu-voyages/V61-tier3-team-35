package game

import "testing"

func TestEvaluateGuess_AllCorrect(t *testing.T) {
	result := EvaluateGuess("apple", "apple")

	if !result.IsCorrect {
		t.Fatal("expected guess to be correct")
	}

	for _, feedback := range result.Feedback {
		if feedback.Status != Correct {
			t.Fatalf("expected all statuses to be correct, got %s",
				feedback.Status)
		}
	}
}

func TestEvaluateGuess_WithDuplicateLetters(t *testing.T) {
	result := EvaluateGuess("alley", "apple")

	expected := []string{
		Correct,
		Present,
		Absent,
		Present,
		Absent,
	}

	for i, expectedStatus := range expected {
		if result.Feedback[i].Status != expectedStatus {
			t.Fatalf("index %d: expected %s, got %s", i,
				expectedStatus, result.Feedback[i].Status)
		}
	}
}

func TestEvaluateGuess_AllAbsent(t *testing.T) {
	result := EvaluateGuess("crony", "apple")

	for _, feedback := range result.Feedback {
		if feedback.Status != Absent {
			t.Fatalf("expected all statuses to be absent, got %s",
				feedback.Status)
		}
	}
}
