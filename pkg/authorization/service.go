package authorization

import (
	"context"
	"fmt"
	"net/http"

	"github.com/proj-go-5/accounts/pkg/accountsio"
	"github.com/proj-go-5/accounts/pkg/jwt"
)

type Service struct {
	jwtService *jwt.Service
}

func NewAuthServie(t *jwt.Service) *Service {
	return &Service{
		jwtService: t,
	}
}

func (s *Service) AdminMiddleware(nextHandler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			accountsio.MakeResponse(w, "unautorized", http.StatusUnauthorized)
			return
		}

		jwtToken, err := s.jwtService.VerifyToken(token)

		if err != nil {
			accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, err := s.jwtService.ExtractClaims(jwtToken)

		if err != nil {
			accountsio.MakeResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "UserId", fmt.Sprintf("%v", claims.ID))
		ctx = context.WithValue(ctx, "UserLogin", claims.Login)

		nextHandler(w, r.WithContext(ctx))
	}
}