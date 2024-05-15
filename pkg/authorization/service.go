package authorization

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/proj-go-5/accounts/pkg/accountsio"
)

type Service struct {
	jwtService *JwtService
}

func NewAuthServie(j *JwtService) *Service {
	return &Service{
		jwtService: j,
	}
}

func (s *Service) AdminMiddleware(nextHandler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			accountsio.MakeResponse(w, "unautorized", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")

		tokenPrefix := "Bearer"

		if len(tokenParts) != 2 || tokenParts[0] != tokenPrefix {
			errorMsg := fmt.Sprintf("The value of the Authorization header must be of the next format: '%v <you_token_value>'", tokenPrefix)
			accountsio.MakeResponse(w, errorMsg, http.StatusUnauthorized)
			return
		}
		token := tokenParts[1]

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
