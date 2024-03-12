package tool

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type TokenManagerer interface {
	Sign(id string, username string) (string, error)
	Unsign(token string) (Claims, error)
}

type Claims struct {
	ID       string `json:"id" `
	Username string `json:"username"`
	jwt.RegisteredClaims
}
type TokenManager struct {
	Token []byte
}

func NewTokenManager(token []byte) *TokenManager {
	return &TokenManager{
		Token: token,
	}
}

func (t *TokenManager) Sign(id, username string) (string, error) {
	claims := &Claims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sined, err := token.SignedString(t.Token)
	if err != nil {
		return "", err
	}
	return sined, nil
}

func (t *TokenManager) Unsign(token string) (Claims, error) {
	log.Println("TOKEN INTO TOKEN MANAGER ", token)
	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return t.Token, nil
	})
	if err != nil {
		return Claims{}, err
	}
	switch {
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return Claims{}, fmt.Errorf("Token is expired")
	}
	cl, ok := tk.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, fmt.Errorf("Can't map claims to token")
	}
	claim := Claims{
		ID:       cl["id"].(string),
		Username: cl["username"].(string),
	}
	return claim, nil

}
