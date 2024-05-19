package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/proj-go-5/accounts/internal/api"
	"github.com/proj-go-5/accounts/internal/api/dto"
	"github.com/proj-go-5/accounts/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestService struct {
}

func (s *TestService) initServices() error {
	return nil
}

func (s *TestService) initRouter() error {
	return nil
}

func (s *TestService) SetUp() error {
	err := s.initServices()
	if err != nil {
		return err
	}

	err = s.initRouter()
	if err != nil {
		return err
	}
	return nil
}

func (s *TestService) CleanUp() {
}

func TestCreateAdminOk(t *testing.T) {
	testService := TestService{}
	err := testService.SetUp()
	defer testService.CleanUp()

	if err != nil {
		return
	}

	app, err := api.NewApp()
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		app.Start()
	}()

	assert.Eventually(t, func() bool {
		_, err := http.Get("http://localhost:8081")
		return err == nil
	}, time.Second, 10*time.Millisecond)

	adminLogin := "test"
	adminPassword := "test"

	app.Service.Admin.Save(&entities.AdminWithPassword{
		Login:    adminLogin,
		Password: adminPassword,
	})

	requestData := &dto.LoginRequest{
		Login:    adminLogin,
		Password: adminPassword,
	}
	data, err := json.Marshal(requestData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/api/v1/auth/admins/login/", bytes.NewReader([]byte(data)))
	if err != nil {
		log.Fatal(err)
		return
	}
	cli := http.DefaultClient
	resp, err := cli.Do(req)

	require.NoError(t, err)

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	app.Stop()
}

func TestCreateAdminFail(t *testing.T) {
	testService := TestService{}
	err := testService.SetUp()
	defer testService.CleanUp()

	if err != nil {
		log.Fatal(err)
		return
	}

	app, err := api.NewApp()
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		app.Start()
	}()

	assert.Eventually(t, func() bool {
		_, err := http.Get("http://localhost:8081")
		return err == nil
	}, time.Second, 10*time.Millisecond)

	adminLogin := "test"
	adminPassword := "test"

	requestData := &dto.LoginRequest{
		Login:    adminLogin,
		Password: adminPassword,
	}
	data, err := json.Marshal(requestData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/api/v1/auth/admins/login/", bytes.NewReader([]byte(data)))
	if err != nil {
		log.Fatal(err)
		return
	}
	cli := http.DefaultClient
	resp, err := cli.Do(req)
	require.NoError(t, err)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.True(t, strings.Contains(string(bodyBytes), "wrong login or password"))

	app.Stop()
}
