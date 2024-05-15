package authorization

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type AdminClaims struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type JwtService struct {
	secret     string
	expiration int
}

func NewJwtService(secret string, expiration int) *JwtService {
	return &JwtService{secret: secret, expiration: expiration}
}

func (s *JwtService) Generate(u *AdminClaims) (string, error) {
	jsonUser, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": string(jsonUser),
		"exp":   time.Now().Add(time.Hour * time.Duration(s.expiration)).Unix(),
	})

	signedToken, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *JwtService) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (s *JwtService) ExtractClaims(token *jwt.Token) (*AdminClaims, error) {
	var userClaim AdminClaims
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
