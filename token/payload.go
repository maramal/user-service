package token

import (
	"errors"
	"time"
)

var (
	ErrInvalidToken = errors.New("token inválido")
	ErrExpiredToken = errors.New("token expirado")
)

type Payload struct {
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	UserType     string    `json:"user_type"`
	ProfileImage string    `json:"profile_image"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiredAt    time.Time `json:"expired_at"`
}

// Crea un nuevo token para un usuario y duración específicos
func NewPayload(fname, lname, email, utype, pimage string, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		FirstName:    fname,
		LastName:     lname,
		Email:        email,
		UserType:     utype,
		ProfileImage: pimage,
		IssuedAt:     time.Now(),
		ExpiredAt:    time.Now().Add(duration),
	}

	return payload, nil
}

// Verifica si el token es válido
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
