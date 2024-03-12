package tool

import (
	"golang.org/x/crypto/bcrypt"
)

type Crypter interface {
	GeneratePassword(string) (string, error)
	CompareHashWithPassowrd(string, string) bool
}
type Bcrypt struct {
}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{}
}

func (b *Bcrypt) GeneratePassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (b *Bcrypt) CompareHashWithPassowrd(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
