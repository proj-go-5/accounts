package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/proj-go-5/accounts/internal/services"
	"github.com/proj-go-5/accounts/pkg/accountsio"
)

type Auth struct {
	tokenService *services.Token
}

func NewAuthServie(t *services.Token) *Auth {
	return &Auth{
		tokenService: t,
	}
}

func (a *Auth) Check(nextHandler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			accountsio.MakeResponse(w, "unautorized", http.StatusUnauthorized)
			return
		}

		jwtToken, err := a.tokenService.VerifyToken(token)

		if err != nil {
			accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, err := a.tokenService.ExtractClaims(jwtToken)

		if err != nil {
			accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "UserId", fmt.Sprintf("%v", claims.ID))
		ctx = context.WithValue(ctx, "UserLogin", claims.Login)

		nextHandler(w, r.WithContext(ctx))
	}
}
