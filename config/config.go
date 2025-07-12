package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App
		HTTP
	}

	App struct {
		Name                  string   `env:"APP_NAME"`
		Version               string   `env:"APP_VERSION"`
		MaxNumTasks           int      `env:"APP_MAX_NUM_TASKS"`
		MaxNumFiles           int      `env:"APP_MAX_NUM_FILES"`
		AllowedFileExtensions []string `env:"APP_ALLOWED_FILE_EXTENSIONS"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT"`
	}
)

var Cfg *Config

func NewConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("unable to load .env file: %w", err)
	}

	Cfg = &Config{}
	if err := env.Parse(Cfg); err != nil {
		return fmt.Errorf("error parce env: %w", err)
	}

	return nil
}
