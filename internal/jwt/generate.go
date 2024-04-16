package jwt

import (
	"fmt"
	"spsu-chat/pkg/clock"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (j *JWT) GeneratePair(userID int64) (TokenPair, error) {
	var tokenPair = TokenPair{}

	var err error
	tokenPair.AccessToken, err = j.GenerateAccessToken(userID)
	if err != nil {
		return tokenPair, err
	}
	tokenPair.RefreshToken, err = j.GenerateRefreshToken(userID)
	if err != nil {
		return tokenPair, err
	}

	return tokenPair, nil
}

func (j *JWT) GenerateAccessToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(clock.Now().Add(j.accessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(clock.Now()),
		Issuer:    j.issuer,
		Subject:   fmt.Sprintf("%d", userID),
	})

	signedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWT) GenerateRefreshToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(clock.Now().Add(j.refreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(clock.Now()),
		Issuer:    j.issuer,
		Subject:   fmt.Sprintf("%d", userID),
	})

	signedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
