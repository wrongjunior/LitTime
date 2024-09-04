package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DefaultReadingSpeed int    `mapstructure:"default_reading_speed"`
	DefaultWorkers      int    `mapstructure:"default_workers"`
	OutputFile          string `mapstructure:"output_file"`
}

// LoadConfig загружает конфигурацию из файла config.yaml или использует значения по умолчанию.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")                                   // Имя конфигурационного файла (без расширения)
	viper.SetConfigType("yaml")                                     // Тип файла конфигурации
	viper.AddConfigPath(".")                                        // Текущая директория
	viper.AddConfigPath("/Users/daniilsolovey/Program/go/LitTime/") // Ваша конкретная директория

	// Установим значения по умолчанию
	viper.SetDefault("default_reading_speed", 180)
	viper.SetDefault("default_workers", 4)
	viper.SetDefault("output_file", "littime_results.json")

	// Попытаемся прочитать конфигурацию из файла
	err := viper.ReadInConfig()
	if err != nil {
		// Если конфиг не найден, это не ошибка, будем использовать значения по умолчанию
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Конфигурационный файл не найден, используются значения по умолчанию")
		} else {
			// Возвращаем ошибку, если произошла другая ошибка при чтении файла
			return nil, fmt.Errorf("не удалось прочитать конфигурационный файл: %w", err)
		}
	}

	// Маппим конфигурацию на структуру Config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("не удалось распарсить конфигурацию: %w", err)
	}

	return &config, nil
}
