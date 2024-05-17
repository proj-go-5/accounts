package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/proj-go-5/accounts/internal/api"
	"github.com/proj-go-5/accounts/internal/api/dto"
	"github.com/proj-go-5/accounts/internal/entities"
	store "github.com/proj-go-5/accounts/internal/repositories"
	"github.com/proj-go-5/accounts/internal/services"
	"github.com/proj-go-5/accounts/pkg/authorization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestService struct {
	Router *mux.Router

	AdminService *services.Admin
	HashService  *services.Hash
	EnvService   *services.Env
	CacheService *services.Cache
	JwtServcie   *authorization.JwtService
	AppServcie   *services.AppService

	EnvFilePath string
}

func (s *TestService) initServices() error {
	filename := ".env"
	s.EnvFilePath = filename

	os.WriteFile(filename, []byte(""), 0644)

	envService, err := services.NewEnvService(filename)
	if err != nil {
		return err
	}
	s.EnvService = envService

	hashService := services.NewHashService()
	s.HashService = hashService
	adminService := services.NewAdminService(store.NewUserMemoryRepository(), hashService)
	s.AdminService = adminService
	cacheService := services.NewCacheService(store.NewMemoryCacheRepository())
	s.CacheService = cacheService

	jwtSecret := envService.Get("JWT_SECRET", "secret")
	jwtExpiration, err := strconv.Atoi(envService.Get("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		return err
	}

	jwtService := authorization.NewJwtService(jwtSecret, jwtExpiration)
	s.JwtServcie = jwtService

	appService := &services.AppService{
		Admin: adminService,
		Jwt:   jwtService,
		Cache: cacheService,
		Auth: services.NewAuthService(
			adminService, cacheService, jwtService, hashService,
		),
	}
	s.AppServcie = appService
	return nil
}

func (s *TestService) initRouter() error {
	a := api.New(s.AppServcie)

	router, err := a.CreateRouter()
	if err != nil {
		return err
	}

	s.Router = router

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
	defer os.Remove(s.EnvFilePath)
}

func TestCreateAdminOk(t *testing.T) {
	testService := TestService{}
	err := testService.SetUp()
	defer testService.CleanUp()

	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		if err := http.ListenAndServe(":8888", testService.Router); err != nil {
			log.Printf("Server run error: %s", err)
		}
	}()

	adminLogin := "test"
	adminPassword := "test"

	testService.AdminService.Save(&entities.AdminWithPassword{
		Login:    adminLogin,
		Password: adminPassword,
	})

	requestData := &dto.LoginRequest{
		Login:    adminLogin,
		Password: adminPassword,
	}
	data, err := json.Marshal(requestData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8888/api/v1/auth/admins/login/", bytes.NewReader([]byte(data)))
	if err != nil {
		log.Fatal(err)
		return
	}
	cli := http.DefaultClient
	resp, err := cli.Do(req)

	require.NoError(t, err)

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestCreateAdminFail(t *testing.T) {
	testService := TestService{}
	err := testService.SetUp()
	defer testService.CleanUp()

	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		if err := http.ListenAndServe(":8888", testService.Router); err != nil {
			log.Printf("Server run error: %s", err)
		}
	}()

	adminLogin := "test"
	adminPassword := "test"

	requestData := &dto.LoginRequest{
		Login:    adminLogin,
		Password: adminPassword,
	}
	data, err := json.Marshal(requestData)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8888/api/v1/auth/admins/login/", bytes.NewReader([]byte(data)))
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
}
