package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/proj-go-5/accounts/pkg/accountsio"
	"github.com/proj-go-5/accounts/pkg/env"
)

var VerifyTokenUrlPath = "/api/v1/auth/admins/me/"

func Authorize(nextHandler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			accountsio.MakeResponse(w, "unautorized", http.StatusUnauthorized)
			return
		}

		envService, err := env.NewEnvService(".env")
		if err != nil {
			log.Println(err)
			return
		}

		checkTokenUrl := fmt.Sprintf("%v:%v%v",
			envService.Get("ACCOUNTS_SERVER_HOST", "http://localhost"),
			envService.Get("ACCOUNTS_SERVER_PORT", "8000"),
			VerifyTokenUrlPath,
		)

		req, err := http.NewRequest(http.MethodGet, checkTokenUrl, nil)
		if err != nil {
			accountsio.MakeResponse(w, "internal error", http.StatusInternalServerError)
			return
		}

		req.Header.Add("Authorization", token)

		client := &http.Client{}

		resp, err := client.Do(req)

		if err != nil {
			accountsio.MakeResponse(w, "internal error", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			accountsio.MakeResponse(w, resp.Body, resp.StatusCode)
			return
		}

		nextHandler(w, r)
	}
}
