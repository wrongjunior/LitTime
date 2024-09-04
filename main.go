package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"LitTime/cmd"
	"LitTime/config"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "littime",
		Short: "LitTime - Reading Time Estimator",
		Long:  `LitTime is a tool for estimating reading time of text documents.`,
	}

	rootCmd.AddCommand(cmd.NewRunCmd(cfg))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
