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

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the LitTime estimator",
		RunE: func(cmd *cobra.Command, args []string) error {
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

	cmd.MarkFlagRequired("file")

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
