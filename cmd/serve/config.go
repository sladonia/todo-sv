package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

type Mongo struct {
	DSN                    string        `default:"mongodb://127.0.0.1:27017" env:"MONGO_DSN"`
	ToDoDatabaseName       string        `default:"todo" env:"MONGO_PROJECT_DATABASE"`
	ProjectsCollectionName string        `default:"projects" evn:"MONGO_PROJECTS_COLLECTION"`
	ConnectTimeout         time.Duration `default:"3s" env:"MONGO_CONNECT_TIMEOUT"`
}

type Config struct {
	ServiceName     string        `default:"todo-sv" env:"SERVICE_NAME"`
	LogLevel        string        `default:"debug" env:"LOG_LEVEL"`
	Port            string        `default:"8080" env:"PORT"`
	ShutdownTimeout time.Duration `default:"5s" env:"SHUTDOWN_TIMEOUT"`
	Mongo           Mongo
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
