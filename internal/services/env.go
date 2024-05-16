package services

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	variablesFilePath string
}

func NewEnvService(variablesFilePath string) (*Env, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	return &Env{variablesFilePath: variablesFilePath}, nil
}

func (s *Env) Get(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		value = defaultValue
	}

	return value
}
