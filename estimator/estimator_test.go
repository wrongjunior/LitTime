package estimator

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestCountSyllables(t *testing.T) {
	tests := []struct {
		word     string
		expected int
	}{
		{"кот", 1},
		{"собакен", 3},
		{"преимущество", 5},
		{"cat", 1},
		{"elephant", 3},
		{"beautiful", 4},
		{"университет", 5},
		{"a", 1},
		{"ь", 1},
		{"rhythm", 1},
		{"apple", 2},
		{"simple", 2},
		{"complicated", 4},
		{"здравствуйте", 3},
		{"воздухоплавание", 6},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("word=%s", test.word), func(t *testing.T) {
			result := CountSyllables(test.word)
			if result != test.expected {
				t.Errorf("CountSyllables(%s) = %d; want %d", test.word, result, test.expected)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		text          string
		expectedCount int
		expectedWords []string
	}{
		{
			"Это простой тест. It contains English words too.",
			8,
			[]string{"Это", "простой", "тест", "It", "contains", "English", "words", "too"},
		},
		{
			"One-two three,four five!",
			4,
			[]string{"One-two", "three", "four", "five"},
		},
		{
			"",
			0,
			[]string{},
		},
		{
			"123 456 789",
			3,
			[]string{"123", "456", "789"},
		},
		{
			"Hyphenated-word non-hyphenated word",
			3,
			[]string{"Hyphenated-word", "non-hyphenated", "word"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("text='%s'", test.text), func(t *testing.T) {
			count, words := CountWords(test.text)
			if count != test.expectedCount {
				t.Errorf("CountWords() returned %d words; want %d", count, test.expectedCount)
			}
			if len(words) != len(test.expectedWords) {
				t.Errorf("CountWords() returned %d words; want %d", len(words), len(test.expectedWords))
			}
			for i, word := range words {
				if i < len(test.expectedWords) && word != test.expectedWords[i] {
					t.Errorf("Word at index %d is %s; want %s", i, word, test.expectedWords[i])
				}
			}
		})
	}
}

func TestCountSentences(t *testing.T) {
	tests := []struct {
		text     string
		expected int
	}{
		{"Это первое предложение. Это второе! А это третье?", 3},
		{"Одно. Два.. Три... Четыре!", 4},
		{"Просто текст", 1},
		{"", 0},
		{"Hello! How are you? I'm fine. Thanks!", 4},
		{"Эллипсис... Это интересно.", 2},
		{"Вопрос? Ответ! Ещё вопрос?! Крик!!!", 4},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("text='%s'", test.text), func(t *testing.T) {
			count := CountSentences(test.text)
			if count != test.expected {
				t.Errorf("CountSentences(%q) = %d; want %d", test.text, count, test.expected)
			}
		})
	}
}

func TestFleschKincaidIndex(t *testing.T) {
	tests := []struct {
		wordsCount     float64
		sentencesCount float64
		syllablesCount float64
		expected       float64
	}{
		{100, 10, 150, 69.8},
		{200, 20, 300, 69.8},
		{0, 0, 0, 0},
		{1, 1, 1, 100},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("w=%.0f,s=%.0f,sy=%.0f", test.wordsCount, test.sentencesCount, test.syllablesCount), func(t *testing.T) {
			result := FleschKincaidIndex(test.wordsCount, test.sentencesCount, test.syllablesCount)
			if math.Abs(result-test.expected) > 0.1 {
				t.Errorf("FleschKincaidIndex(%.0f, %.0f, %.0f) = %.2f; want %.2f",
					test.wordsCount, test.sentencesCount, test.syllablesCount, result, test.expected)
			}
		})
	}
}

func TestEstimateReadingTimeParallel(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		readingSpeed float64
		hasVisuals   bool
		workerCount  int
		expectError  bool
	}{
		{
			name:         "Normal text",
			text:         "Это тестовый текст. Он содержит несколько предложений. And some English words too.",
			readingSpeed: 250,
			hasVisuals:   false,
			workerCount:  2,
			expectError:  false,
		},
		{
			name:         "Short text with visuals",
			text:         "Короткий текст.",
			readingSpeed: 200,
			hasVisuals:   true,
			workerCount:  1,
			expectError:  false,
		},
		{
			name:         "Empty text",
			text:         "",
			readingSpeed: 250,
			hasVisuals:   false,
			workerCount:  4,
			expectError:  true,
		},
		{
			name:         "Complex text",
			text:         "Это сложный текст с длинными предложениями и редкими словами. Он содержит много слогов и может быть трудным для чтения. Мы используем его для проверки алгоритма.",
			readingSpeed: 200,
			hasVisuals:   false,
			workerCount:  3,
			expectError:  false,
		},
		{
			name:         "Single word, single sentence",
			text:         "Привет.",
			readingSpeed: 250,
			hasVisuals:   false,
			workerCount:  1,
			expectError:  false,
		},
		{
			name:         "Text with non-word characters",
			text:         "!!! ?? .. --",
			readingSpeed: 250,
			hasVisuals:   false,
			workerCount:  1,
			expectError:  true,
		},
		{
			name:         "Long text with complex structure",
			text:         strings.Repeat("Это длинный текст с многими словами и предложениями. ", 1000),
			readingSpeed: 200,
			hasVisuals:   false,
			workerCount:  8,
			expectError:  false,
		},
		{
			name:         "Text with visuals",
			text:         "Текст с визуальными элементами.",
			readingSpeed: 250,
			hasVisuals:   true,
			workerCount:  2,
			expectError:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := EstimateReadingTimeParallel(test.text, test.readingSpeed, test.hasVisuals, test.workerCount)

			if test.expectError {
				if err == nil {
					t.Error("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("EstimateReadingTimeParallel() returned an unexpected error: %v", err)
				}
				if result.ReadingTime < 0 {
					t.Errorf("EstimateReadingTimeParallel() returned invalid reading time: %f", result.ReadingTime)
				}
				expectedWordCount := len(wordRegex.FindAllString(test.text, -1))
				if result.WordCount != expectedWordCount {
					t.Errorf("Word count mismatch. Got %d, want %d", result.WordCount, expectedWordCount)
				}
				expectedSentences := CountSentences(test.text)
				if result.SentenceCount != expectedSentences {
					t.Errorf("Sentence count mismatch. Got %d, want %d", result.SentenceCount, expectedSentences)
				}
				if result.FleschKincaidIndex < 0 || result.FleschKincaidIndex > 120 {
					t.Errorf("Invalid Flesch-Kincaid Index: %f", result.FleschKincaidIndex)
				}
			}
		})
	}
}
