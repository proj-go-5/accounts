package services

import (
	"golang.org/x/crypto/bcrypt"
)

var workFactor = 14

type Hash struct {
}

func NewHashService() *Hash {
	return &Hash{}
}

func (h *Hash) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), workFactor)
	return string(bytes), err
}

func (h *Hash) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
