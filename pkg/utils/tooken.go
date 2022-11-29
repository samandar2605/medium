package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/post/config"
)

const SecretKey = "%r[(DfWOy2y~9bZ"

type TokenParams struct {
	UserID          int
	UserType        string
	FirstName       string
	LastName        string
	Username        string
	Email           string
	ProfileImageUrl string
	Duration        time.Duration
}

// CreateToken creates a new token
func CreateToken(cfg *config.Config, params *TokenParams) (string, *Payload, error) {
	payload, err := NewPayload(params)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(SecretKey))
	return token, payload, err
}

func VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
