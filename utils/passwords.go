package utils

import (
	"golang.org/x/crypto/bcrypt"
)

/**
 * Generate a hash from a password
 *
 * @param password string The password to hash
 * @return string The hash
 */
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/**
 * Verifica si una contraseña coincide con un hash
 *
 * @param password string "La contraseña a verificar"
 * @param hash string "El hash a verificar"
 * @return error "El error"
 */
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
