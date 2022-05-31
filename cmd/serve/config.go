package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName     string        `default:"todo-sv" env:"SERVICE_NAME"`
	LogLevel        string        `default:"debug" env:"LOG_LEVEL"`
	Port            string        `default:"8080" env:"PORT"`
	ShutdownTimeout time.Duration `default:"5s" env:"SHUTDOWN_TIMEOUT"`
}

func mustLoadConfig() Config {
	var config Config

	_ = godotenv.Load()

	err := configor.Load(&config)
	if err != nil {
		panic(fmt.Sprintf("failed to load config. err: %s", err.Error()))
	}

	return config
}
