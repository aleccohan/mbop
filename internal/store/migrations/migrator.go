package migrations

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/redhatinsights/mbop/internal/config"

	// this is the file:// driver for go-migrate.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(driver database.Driver) error {
	m, err := migrate.NewWithDatabaseInstance("file://migrations", config.Get().DatabaseName, driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
