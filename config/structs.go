package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Log      LogConfig      `mapstructure:"log"`
	Database DatabaseConfig `mapstructure:"database"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Dir    string `mapstructure:"dir"`
	MaxAge int    `mapstructure:"max_age"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

func mapToStructs() {
	loadEnvToStruct("app.name", "APP_NAME")
	loadEnvToStruct("app.env", "APP_ENV")
	loadEnvToStruct("app.port", "APP_PORT")
	loadEnvToStruct("log.level", "LOG_LEVEL")
	loadEnvToStruct("log.dir", "LOG_DIR")
	loadEnvToStruct("log.max_age", "LOG_MAX_AGE")
	loadEnvToStruct("database.path", "DB_PATH")
}

func loadEnvToStruct(key string, env string) {
	if err := viper.BindEnv(key, env); err != nil {
		log.Fatalf("failed to map environment variable into struct field: %v", err)
	}
}
