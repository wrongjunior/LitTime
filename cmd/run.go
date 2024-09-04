package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"LitTime/config"
	"LitTime/estimator"
	"LitTime/ui"
)

func NewRunCmd(cfg *config.Config) *cobra.Command {
	var filePath string
	var readingSpeed int
	var hasVisuals bool
	var workers int
	var interactive bool

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the LitTime estimator",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Если интерактивный режим включен, запускаем интерфейс через bubbletea
			if interactive {
				userInputs, err := ui.RunInteractive(cfg)
				if err != nil {
					return err
				}
				filePath = userInputs.FilePath
				readingSpeed = userInputs.ReadingSpeed
				hasVisuals = userInputs.HasVisuals
				workers = userInputs.Workers
			}

			// Проверяем, если интерактивный режим выключен, то файл должен быть указан через флаг
			if !interactive && filePath == "" {
				return fmt.Errorf("required flag(s) \"file\" not set")
			}

			// Проверяем, был ли передан валидный путь к файлу
			if filePath == "" {
				return fmt.Errorf("file path cannot be empty")
			}

			// Запуск оценки времени чтения
			result, err := runEstimator(filePath, readingSpeed, hasVisuals, workers)
			if err != nil {
				return err
			}

			fmt.Printf("Saving result to: %s\n", cfg.OutputFile)
			if err := saveResult(result, cfg.OutputFile); err != nil {
				return fmt.Errorf("failed to save result: %w", err)
			}

			return ui.RunUI(result)
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the text file")
	cmd.Flags().IntVarP(&readingSpeed, "speed", "s", cfg.DefaultReadingSpeed, "Reading speed in words per minute")
	cmd.Flags().BoolVarP(&hasVisuals, "visuals", "v", false, "Set to true if the text contains visual elements")
	cmd.Flags().IntVarP(&workers, "workers", "w", cfg.DefaultWorkers, "Number of worker goroutines")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Enable interactive mode for setting options")

	return cmd
}

func runEstimator(filePath string, readingSpeed int, hasVisuals bool, workers int) (*estimator.Result, error) {
	// Чтение текста из файла
	text, err := estimator.ReadTextFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Оценка времени чтения
	result, err := estimator.EstimateReadingTimeParallel(text, float64(readingSpeed), hasVisuals, workers)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate reading time: %w", err)
	}

	return &result, nil
}

func saveResult(result *estimator.Result, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}
