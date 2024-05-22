package services

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvServiceOk(t *testing.T) {
	varName := "TEST_VAR"
	varValue := "test"

	content := fmt.Sprintf("%v=%v", varName, varValue)

	filename := "env.test"

	os.WriteFile(filename, []byte(content), 0644)

	defer os.Remove(filename)

	envService, _ := NewEnvService(filename)

	envVarValue := envService.Get(varName, "")

	assert.Equal(t, envVarValue, varValue)

	defaultValue := "default"
	noneExistedVarName := "NONE_EXISTED_VAR"

	envVarValue = envService.Get(noneExistedVarName, defaultValue)

	assert.Equal(t, envVarValue, defaultValue)
}

func TestEnvServiceFail(t *testing.T) {
	noneExistedFileName := "nonexisted.env"

	_, err := NewEnvService(noneExistedFileName)
	assert.Equal(t, err.Error(), fmt.Sprintf("open %v: no such file or directory", noneExistedFileName))
}
