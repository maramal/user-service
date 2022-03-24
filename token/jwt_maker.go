package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

/** Crea un nuevo JWTMaker
 *
 * @param secretKey string "Clave secreta"
 * @return *JWTMaker "Instancia de JWTMaker"
 * @return error "Error"
 */
func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("el tamaño de la clave secreta debe ser de al menos %d caracteres", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

/** Crea un nuevo token para un usuario y duración específicos
 *
 * @param fname string "Nombre"
 * @param lname string "Apellido"
 * @param email string "Email del usuario"
 * @param utype string "Tipo de usuario"
 * @param pimage string "Imagen de perfil"
 * @param duration time.Duration "Duración del token"
 * @return string "Token"
 * @return *Payload "Payload del token"
 * @return error "Error"
 */
func (maker *JWTMaker) CreateToken(fname, lname, email, utype, pimage string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(fname, lname, email, utype, pimage, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

/** Verifica que el token sea válido o no
 *
 * @param token string "Token"
 * @return *Payload "Payload del token"
 * @return error "Error"
 */
func (maker *JWTMaker) Valid(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
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
