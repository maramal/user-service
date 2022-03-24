package token

import (
	"time"
)

// Maker es una interface para administrar tokens
type IMaker interface {
	// Crea un nuevo token para un usuario y duración específicos
	CreateToken(fname, lname, email, utype, pimage string, duration time.Duration) (string, *Payload, error)

	// Verifica si el token es válido o no
	Valid(token string) (*Payload, error)
}
