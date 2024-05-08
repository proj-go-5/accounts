package services

import (
	"encoding/json"
	"errors"
	"github.com/proj-go-5/accounts/internal/entities"
	"time"

	"github.com/golang-jwt/jwt"
)

var tokenSecret = "replace_me_by_env_var!!!"

type Token struct {
}

func NewTokenService() *Token {
	return &Token{}
}

func (t *Token) Generate(u *entities.User) (string, error) {
	jsonUser, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": string(jsonUser),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (t *Token) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (t *Token) ExtractClaims(token *jwt.Token) (*entities.UserClaims, error) {
	var userClaim entities.UserClaims
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return &userClaim, errors.New("invalid claims format")
	}

	userJson, ok := claims["user"].(string)
	if !ok {
		return &userClaim, errors.New("subject claim not found")
	}

	if err := json.Unmarshal([]byte(userJson), &userClaim); err != nil {
		return &userClaim, err
	}

	return &userClaim, nil
}
