package config

import (
	"errors"
	"log"

	"github.com/adnanahmady/go-rest-api-blog/pkg/app"
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
	root := app.GetRootPath()
	if err := godotenv.Load(
		root + "/.env.testing",
		root + "/.env",
	); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
}

func readConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(app.GetRootPath())

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
