package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	UsersTable = "users"
)

func GetPgError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}

	return nil
}

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Config struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	Port     uint   `yaml:"port" env:"POSTGRES_PORT"`
	Database string `yaml:"database" env:"POSTGRES_DATABASE"`
	User     string `yaml:"user" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
}

type PostgresqlRepository struct {
	DB DB
	db *sqlx.DB
}

func New(config Config) (PostgresqlRepository, error) {
	db, err := sqlx.Open("pgx", formatConnectionUrl(config))
	if err != nil {
		return PostgresqlRepository{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return PostgresqlRepository{}, fmt.Errorf("failed to ping database: %w", err)
	}
	return PostgresqlRepository{
		DB: db,
		db: db,
	}, nil
}

func formatConnectionUrl(config Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
}
