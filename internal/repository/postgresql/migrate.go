package postgresql

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(ctx context.Context, config Config) error {
	migrate, err := migrate.New("file://schema", formatConnectionUrl(config))
	if err != nil {
		return err
	}
	return migrate.Up()
}
