package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var (
	russianVowels    = "аеёиоуыэюя"
	englishVowels    = "aeiouy"
	wordRegex        = regexp.MustCompile(`\p{L}+`)
	sentenceEndRegex = regexp.MustCompile(`[.!?]+`)
)

func countSyllables(word string) int {
	word = strings.ToLower(word)
	syllables := 0
	isRussian := false

	// определяем язык слова по первой букве
	if len(word) > 0 && unicode.Is(unicode.Cyrillic, rune(word[0])) {
		isRussian = true
	}

	vowels := englishVowels
	if isRussian {
		vowels = russianVowels
	}

	if strings.ContainsRune(vowels, rune(word[0])) {
		syllables++
	}

	for i := 1; i < len(word); i++ {
		if strings.ContainsRune(vowels, rune(word[i])) && !strings.ContainsRune(vowels, rune(word[i-1])) {
			syllables++
		}
	}

	if isRussian && (strings.HasSuffix(word, "ь") || strings.HasSuffix(word, "й")) {
		syllables--
	}

	return max(syllables, 1)
}

func countWords(text string) (int, []string) {
	words := wordRegex.FindAllString(text, -1)
	return len(words), words
}

func countSentences(text string) int {
	sentences := sentenceEndRegex.Split(text, -1)
	return len(sentences)
}

func fleschKincaidIndex(wordsCount, sentencesCount, syllablesCount float64) float64 {
	return 206.835 - 1.3*(wordsCount/sentencesCount) - 60.1*(syllablesCount/wordsCount)
}

func estimateReadingTime(text string, readingSpeed float64, hasVisuals bool) float64 {
	wordsCount, words := countWords(text)
	sentencesCount := countSentences(text)
	syllablesCount := 0

	for _, word := range words {
		syllablesCount += countSyllables(word)
	}

	fkIndex := fleschKincaidIndex(float64(wordsCount), float64(sentencesCount), float64(syllablesCount))

	adjustedSpeed := readingSpeed
	if fkIndex < 60 {
		adjustedSpeed *= 0.8 // Сложный текст
	}

	readingTime := float64(wordsCount) / adjustedSpeed

	if hasVisuals {
		readingTime *= 1.1
	}

	return math.Round(readingTime*100) / 100
}

func readTextFromFile(filePath string) (string, error) {
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

func main() {
	filePath := "example.txt"

	text, err := readTextFromFile(filePath)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	readingTimeMinutes := estimateReadingTime(text, 250, false) // TODO: вынести скорость чтения в глобальную переменную, а потом вообще сделать конфиг
	fmt.Printf("Примерное время на чтение файла '%s': %.2f минут.\n", filePath, readingTimeMinutes)
}
