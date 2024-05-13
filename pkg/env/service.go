package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Service struct {
	variablesFilePath string
}

func NewEnvService(variablesFilePath string) (*Service, error) {
	err := godotenv.Overload(".env")
	if err != nil {
		return nil, err
	}

	return &Service{variablesFilePath: variablesFilePath}, nil
}

func (s *Service) Get(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		value = defaultValue
	}

	return value
}
