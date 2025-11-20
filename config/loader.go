package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var cfg *Config

func LoadConfig() {
	loadDotEnvFile()
	readConfigFile()
	loadEnvVariables()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}

func GetConfig() *Config {
	if cfg == nil {
		LoadConfig()
	}
	return cfg
}

func loadDotEnvFile() {
	if err := godotenv.Load(getCallerPath() + "/.env"); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
}

func getCallerPath() string {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current path: %v", err)
	}
	return currentPath
}

func readConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(getCallerPath())

	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) {
			log.Fatalf("config file not found: %v", err)
			return
		}
		log.Fatalf("failed to read config file: %v", err)
	}
}

func loadEnvVariables() {
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	mapToStructs()
}
