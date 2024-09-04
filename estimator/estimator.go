package estimator

import (
	"bufio"
	"errors"
	"math"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

var (
	russianVowels    = "аеёиоуыэюя"
	englishVowels    = "aeiouy"
	wordRegex        = regexp.MustCompile(`[\p{L}\p{N}]+`)
	sentenceEndRegex = regexp.MustCompile(`[.!?]+`)
)

// содержит результаты анализа текста
type Result struct {
	ReadingTime        float64
	WordCount          int
	SentenceCount      int
	SyllableCount      int
	FleschKincaidIndex float64
}

// подсчитывает количество слогов в слове
func CountSyllables(word string) int {
	word = strings.ToLower(word)
	syllables := 0
	isRussian := false

	if len(word) > 0 && unicode.Is(unicode.Cyrillic, rune(word[0])) {
		isRussian = true
	}

	vowels := englishVowels
	if isRussian {
		vowels = russianVowels
	}

	for i, char := range word {
		if strings.ContainsRune(vowels, char) {
			if i == 0 || !strings.ContainsRune(vowels, rune(word[i-1])) {
				syllables++
			}
		}
	}

	if isRussian && (strings.HasSuffix(word, "ь") || strings.HasSuffix(word, "й")) {
		syllables = max(syllables-1, 1)
	}

	return max(syllables, 1)
}

// подсчитывает количество слов в тексте
func CountWords(text string) (int, []string) {
	words := wordRegex.FindAllString(text, -1)
	return len(words), words
}

// подсчитывает количество предложений в тексте
func CountSentences(text string) int {
	sentences := sentenceEndRegex.Split(strings.TrimSpace(text), -1)
	count := 0
	for _, s := range sentences {
		if strings.TrimSpace(s) != "" {
			count++
		}
	}
	return count
}

// рассчитывает индекс Флеша-Кинкейда
func FleschKincaidIndex(wordsCount, sentencesCount, syllablesCount float64) float64 {
	if wordsCount == 0 || sentencesCount == 0 {
		return 0
	}
	return 206.835 - 1.3*(wordsCount/sentencesCount) - 60.1*(syllablesCount/wordsCount)
}

// оценивает время чтения текста с использованием параллельной обработки
func EstimateReadingTimeParallel(text string, readingSpeed float64, hasVisuals bool, workerCount int) (Result, error) {
	wordsCount, words := CountWords(text)
	sentencesCount := CountSentences(text)

	if wordsCount == 0 || sentencesCount == 0 {
		return Result{}, errors.New("text is empty or invalid")
	}

	// параллельный подсчет слогов
	syllablesChan := make(chan int)
	var wg sync.WaitGroup

	chunkSize := (len(words) + workerCount - 1) / workerCount
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := min(start+chunkSize, len(words))
			localSum := 0
			for j := start; j < end; j++ {
				localSum += CountSyllables(words[j])
			}
			syllablesChan <- localSum
		}(i * chunkSize)
	}

	go func() {
		wg.Wait()
		close(syllablesChan)
	}()

	syllablesCount := 0
	for count := range syllablesChan {
		syllablesCount += count
	}

	fkIndex := FleschKincaidIndex(float64(wordsCount), float64(sentencesCount), float64(syllablesCount))

	adjustedSpeed := readingSpeed
	if fkIndex < 60 {
		adjustedSpeed *= 0.8 // сложный текст
	}

	readingTime := float64(wordsCount) / adjustedSpeed

	if hasVisuals {
		readingTime *= 1.1
	}

	return Result{
		ReadingTime:        math.Round(readingTime*100) / 100,
		WordCount:          wordsCount,
		SentenceCount:      sentencesCount,
		SyllableCount:      syllablesCount,
		FleschKincaidIndex: fkIndex,
	}, nil
}

// читает текст из файла
func ReadTextFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var text strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text.WriteString(scanner.Text())
		text.WriteString(" ")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return text.String(), nil
}
