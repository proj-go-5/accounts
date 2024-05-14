package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/proj-go-5/accounts/internal/entities"

	"github.com/golang-jwt/jwt"
)

var tokenSecret = "replace_me_by_env_var!!!"

type Token struct {
}

func NewTokenService() *Token {
	return &Token{}
}

func (t *Token) Generate(u *entities.Admin) (string, error) {
	jsonUser, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": string(jsonUser),
		"exp":   fmt.Sprintf("%v", time.Now().Add(time.Second*30).Unix()),
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
		fmt.Println(222)
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (t *Token) ExtractClaims(token *jwt.Token) (result *entities.TokenClaims, error error) {
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return result, errors.New("invalid claims format")
	}

	var admin entities.Admin
	adminJson := fmt.Sprintf("%v", claims["admin"])

	json.Unmarshal([]byte(adminJson), &admin)

	fmt.Println("admin: ", admin, admin.ID, admin.Login)
	res := &entities.TokenClaims{
		Exp:   claims["exp"].(int64),
		Admin: admin}

	fmt.Println("res : ", res)
	return res, nil
}
