package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidAlgorithm = errors.New("invalid jwt encryption algorithm")
	ErrInvalidClaims    = errors.New("invalid jwt claims")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidToken     = errors.New("invalid token")
)

type Claims struct {
	ExpiredAt time.Time
	IssuedAt  time.Time
	Issuer    string
	Subject   int64
}

func (j *JWT) ValidateToken(token string) (Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidAlgorithm
		}

		return []byte(j.secret), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return Claims{}, ErrTokenExpired
		}

		return Claims{}, ErrInvalidToken
	}
	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return Claims{}, ErrInvalidClaims
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return Claims{}, ErrInvalidClaims
	}

	parsedClaims := Claims{
		ExpiredAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
		Issuer:    claims.Issuer,
		Subject:   userID,
	}

	return parsedClaims, nil
}
