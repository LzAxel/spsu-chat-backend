package service

import (
	"context"
	"errors"
	"fmt"

	"spsu-chat/internal/apperror"
	"spsu-chat/internal/jwt"
	"spsu-chat/internal/models"
	"spsu-chat/internal/repository"
	"spsu-chat/pkg/clock"
	"spsu-chat/pkg/hash"
)

type AuthorizationSerive struct {
	jwt      *jwt.JWT
	userRepo repository.User
}

func NewAuthorizationSerive(jwt *jwt.JWT, userRepo repository.User) *AuthorizationSerive {
	return &AuthorizationSerive{
		jwt:      jwt,
		userRepo: userRepo,
	}
}

func (a *AuthorizationSerive) RefreshTokens(ctx context.Context, refreshToken string) (jwt.TokenPair, error) {
	claims, err := a.jwt.ValidateToken(refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrInvalidToken) || errors.Is(err, jwt.ErrInvalidClaims):
			return jwt.TokenPair{}, jwt.ErrInvalidToken
		case errors.Is(err, jwt.ErrTokenExpired):
			return jwt.TokenPair{}, jwt.ErrTokenExpired
		}
		return jwt.TokenPair{}, fmt.Errorf("Authorization.RefreshTokens: %w", err)
	}

	tokenPair, err := a.jwt.GeneratePair(claims.Subject)
	if err != nil {
		return jwt.TokenPair{}, fmt.Errorf("Authorization.RefreshTokens: %w", err)
	}

	return tokenPair, err
}

func (a *AuthorizationSerive) Login(ctx context.Context, input models.LoginUserInput) (jwt.TokenPair, error) {
	user, err := a.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		if errors.As(err, &apperror.DBError{}) {
			dbErr := err.(apperror.DBError)
			if errors.Is(dbErr.Err, apperror.ErrNotFound) {
				return jwt.TokenPair{}, models.ErrInvalidCredentials
			}
		}

		return jwt.TokenPair{}, err
	}

	if err := hash.Compare(user.PasswordHash, input.Password); err != nil {
		return jwt.TokenPair{}, models.ErrInvalidCredentials
	}

	tokenPair, err := a.jwt.GeneratePair(user.ID)
	if err != nil {
		return jwt.TokenPair{}, err
	}

	return tokenPair, err
}
func (a *AuthorizationSerive) Register(ctx context.Context, input models.CreateUserInput) error {
	passwordHash, err := hash.Hash(input.Password)
	if err != nil {
		return err
	}

	if input.DisplayName == nil {
		input.DisplayName = &input.Username
	}

	dto := models.CreateUserRecord{
		Username:     input.Username,
		DisplayName:  *input.DisplayName,
		PasswordHash: passwordHash,
		Type:         models.UserTypeUser,
		CreatedAt:    clock.Now(),
	}

	return a.userRepo.Create(ctx, dto)
}
