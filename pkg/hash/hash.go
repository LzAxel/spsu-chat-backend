package hash

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	saltLength = 8
)

var (
	ErrInvalidPassword = fmt.Errorf("invalid password")
)

func Compare(hashedPassword []byte, password string) error {
	salt := hashedPassword[:saltLength]

	pass := []byte(password)
	pass = append(pass, salt...)

	err := bcrypt.CompareHashAndPassword(hashedPassword[saltLength:], []byte(pass))
	if err != nil {
		return ErrInvalidPassword
	}

	return nil
}

func Hash(s string) ([]byte, error) {
	salt, err := generateSalt()
	if err != nil {
		return []byte{}, fmt.Errorf("gen salt for hashing password: %w", err)
	}

	password := []byte(s)
	password = append(password, salt...)

	passwordHash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, fmt.Errorf("hashing password: %w", err)
	}

	result := []byte(salt)
	result = append(result, passwordHash...)

	return result, nil
}

func generateSalt() ([]byte, error) {
	var salt = make([]byte, saltLength)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("error generating salt: %w", err)
	}

	return salt, nil
}
