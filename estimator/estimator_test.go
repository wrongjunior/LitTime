package estimator

import (
	"testing"
)

func TestCountSyllables(t *testing.T) {
	tests := []struct {
		word     string
		expected int
	}{
		{"кот", 1},
		{"собака", 3},
		{"преимущество", 5},
		{"cat", 1},
		{"elephant", 3},
	}

	for _, test := range tests {
		result := CountSyllables(test.word)
		if result != test.expected {
			t.Errorf("CountSyllables(%s) = %d; want %d", test.word, result, test.expected)
		}
	}
}

func TestCountWords(t *testing.T) {
	text := "Это простой тест. It contains English words too."
	count, words := CountWords(text)
	expectedCount := 7
	if count != expectedCount {
		t.Errorf("CountWords() returned %d words; want %d", count, expectedCount)
	}
	if len(words) != expectedCount {
		t.Errorf("CountWords() returned %d words; want %d", len(words), expectedCount)
	}
}

func TestCountSentences(t *testing.T) {
	text := "Это первое предложение. Это второе! А это третье?"
	count := CountSentences(text)
	expected := 3
	if count != expected {
		t.Errorf("CountSentences() returned %d sentences; want %d", count, expected)
	}
}

func TestEstimateReadingTimeParallel(t *testing.T) {
	text := "Это тестовый текст. Он содержит несколько предложений. And some English words too."
	result, err := EstimateReadingTimeParallel(text, 250, false, 2)
	if err != nil {
		t.Fatalf("EstimateReadingTimeParallel() returned an error: %v", err)
	}
	if result.ReadingTime <= 0 {
		t.Errorf("EstimateReadingTimeParallel() returned invalid reading time: %f", result.ReadingTime)
	}
}
