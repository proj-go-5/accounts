package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/proj-go-5/accounts/internal/entities"

	"github.com/golang-jwt/jwt"
)

type Token struct {
	secret     string
	expiration int
}

func NewTokenService(secret string, expiration int) *Token {
	return &Token{secret: secret, expiration: expiration}
}

func (t *Token) Generate(u *entities.Admin) (string, error) {
	jsonUser, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": string(jsonUser),
		"exp":   time.Now().Add(time.Hour * time.Duration(t.expiration)).Unix(),
	})

	signedToken, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (t *Token) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (t *Token) ExtractClaims(token *jwt.Token) (*entities.AdminClaims, error) {
	var userClaim entities.AdminClaims
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return &userClaim, errors.New("invalid claims format")
	}

	userJson, ok := claims["admin"].(string)
	if !ok {
		return &userClaim, errors.New("admins claims not found")
	}

	if err := json.Unmarshal([]byte(userJson), &userClaim); err != nil {
		return &userClaim, err
	}

	return &userClaim, nil
}
