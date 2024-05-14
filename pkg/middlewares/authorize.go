package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/proj-go-5/accounts/internal/services"
	"github.com/proj-go-5/accounts/pkg/accountsio"
)

var VerifyTokenUrlPath = "/api/v1/auth/admins/me/"

func Authorize(nextHandler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			accountsio.MakeResponse(w, "unautorized", http.StatusUnauthorized)
			return
		}

		// envService, err := env.NewEnvService(".env")
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }

		tokenService := services.NewTokenService()

		jwtToken, err := tokenService.VerifyToken(token)

		if err != nil {
			fmt.Println(7777)
			accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokenClaims, err := tokenService.ExtractClaims(jwtToken)

		if err != nil {
			fmt.Println(88888)
			accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		expTime, err := time.Parse(time.RFC3339, tokenClaims.Exp)
		if err != nil {
			fmt.Println(9999)
			accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(expTime.Year())
		fmt.Println(expTime.Month())
		fmt.Println(expTime.Day())

		nextHandler(w, r)
	}
}
