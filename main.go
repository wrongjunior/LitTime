package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	_ "unicode"
)

// подсчитывает количество слогов в слове на основе гласных букв
func countSyllables(word string) int {
	vowels := "аеёиоуыэюя"
	syllables := 0
	word = strings.ToLower(word)
	if strings.ContainsRune(vowels, rune(word[0])) {
		syllables++
	}
	for i := 1; i < len(word); i++ {
		if strings.ContainsRune(vowels, rune(word[i])) && !strings.ContainsRune(vowels, rune(word[i-1])) {
			syllables++
		}
	}
	if strings.HasSuffix(word, "ь") || strings.HasSuffix(word, "й") {
		syllables--
	}
	if syllables < 1 {
		return 1
	}
	return syllables
}

// подсчитывает количество слов в тексте
func countWords(text string) (int, []string) {
	words := regexp.MustCompile(`\b\w+\b`).FindAllString(text, -1)
	return len(words), words
}

// подсчитывает количество предложений в тексте
func countSentences(text string) int {
	sentences := regexp.MustCompile(`[.!?]+`).Split(text, -1)
	if sentences[len(sentences)-1] == "" {
		return len(sentences) - 1
	}
	return len(sentences)
}

// рассчитывает индекс Флеша-Иванова для русского языка
func fleschKincaidIndex(wordsCount, sentencesCount, syllablesCount int) float64 {
	return 206.835 - 1.3*float64(wordsCount)/float64(sentencesCount) - 60.1*float64(syllablesCount)/float64(wordsCount)
}

// оценивает время на чтение текста с учетом сложности и структурных элементов
func estimateReadingTime(text string, readingSpeed float64, hasVisuals bool) float64 {
	wordsCount, words := countWords(text)
	sentencesCount := countSentences(text)
	syllablesCount := 0

	for _, word := range words {
		syllablesCount += countSyllables(word)
	}

	fkIndex := fleschKincaidIndex(wordsCount, sentencesCount, syllablesCount)

	// корректировка скорости чтения в зависимости от сложности текста
	adjustedSpeed := readingSpeed
	if fkIndex < 60 {
		adjustedSpeed *= 0.8 // Сложный текст
	}

	// оценка основного времени на чтение
	readingTime := float64(wordsCount) / adjustedSpeed

	// коррекция на структурные элементы
	if hasVisuals {
		readingTime *= 1.1
	}

	return math.Round(readingTime*100) / 100
}

// считывает содержимое текстового файла и возвращает его как строку.
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
	filePath := "example.txt" // замените на путь к вашему файлу

	text, err := readTextFromFile(filePath)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	readingTime := estimateReadingTime(text, 250, true)
	fmt.Printf("Примерное время на чтение файла '%s': %.2f минут.\n", filePath, readingTime)
}
