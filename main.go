package main

import (
	"flag"
	"fmt"
	"log"
	_ "os"
	"runtime"
	"time"

	"LitTime/estimator"
)

func main() {
	// Настройка флагов командной строки
	filePath := flag.String("file", "", "Path to the text file")
	readingSpeed := flag.Float64("speed", 180, "Reading speed in words per minute")
	hasVisuals := flag.Bool("visuals", false, "Set to true if the text contains visual elements")
	workerCount := flag.Int("workers", runtime.NumCPU(), "Number of worker goroutines")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("Please provide a file path using the -file flag")
	}

	// Настройка логирования
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Обработка паники
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Чтение файла
	start := time.Now()
	text, err := estimator.ReadTextFromFile(*filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Оценка времени чтения
	result, err := estimator.EstimateReadingTimeParallel(text, *readingSpeed, *hasVisuals, *workerCount)
	if err != nil {
		log.Fatalf("Error estimating reading time: %v", err)
	}

	duration := time.Since(start)

	// Вывод результатов
	fmt.Printf("File: %s\n", *filePath)
	fmt.Printf("Estimated reading time: %.2f minutes\n", result.ReadingTime)
	fmt.Printf("Word count: %d\n", result.WordCount)
	fmt.Printf("Sentence count: %d\n", result.SentenceCount)
	fmt.Printf("Syllable count: %d\n", result.SyllableCount)
	fmt.Printf("Flesch-Kincaid Index: %.2f\n", result.FleschKincaidIndex)
	fmt.Printf("Processing time: %v\n", duration)
}
