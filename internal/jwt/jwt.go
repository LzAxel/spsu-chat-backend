package jwt

import (
	"time"
)

type Config struct {
	Secret          string        `yaml:"secret" env:"SECRET"`
	AccessTokenTTL  time.Duration `yaml:"accessTokenTTL" env:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL time.Duration `yaml:"refreshTokenTTL" env:"REFRESH_TOKEN_TTL"`
	Issuer          string        `yaml:"issuer" env:"ISSUER"`
}

type JWT struct {
	secret          string
	issuer          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func New(config Config) *JWT {
	return &JWT{
		secret:          config.Secret,
		accessTokenTTL:  config.AccessTokenTTL,
		refreshTokenTTL: config.RefreshTokenTTL,
		issuer:          config.Issuer,
	}
}
