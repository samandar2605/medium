package utils

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID              uuid.UUID `json:"id"`
	UserId          int       `json:"user_id"`
	UserType        string    `json:"user_type"`
	FirstName       string    `json:"first_name"`
	Username        string    `json:"username"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	ProfileImageUrl *string    `json:"profile_image_url"`
	IssuedAt        time.Time `json:"issued_at"`
	ExpiredAt       time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload
func NewPayload(params *TokenParams) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserId:    params.UserID,
		Email:     params.Email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		UserType:  params.UserType,
		Username:  params.Username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(params.Duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
