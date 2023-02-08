package store

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/store/migrations"
)

// GetStore is a function that will return the currently configured store. this
// allows it to be overridden for testing or alternative implementations
var GetStore func() (Store, error)

func SetupStore() error {
	switch config.Get().StoreBackend {
	case "postgres":
		pgStore, err := setupPostgresStore()
		if err != nil {
			return err
		}

		GetStore = func() (Store, error) { return pgStore, nil }
	case "memory":
		GetStore = func() (Store, error) { return &inMemoryStore{}, nil }
	}

	return nil
}

func setupPostgresStore() (*postgresStore, error) {
	c := config.Get()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=prefer",
		c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort, c.DatabaseName)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return nil, err
	}

	err = migrations.Migrate(driver)
	if err != nil {
		return nil, err
	}

	return &postgresStore{db: db}, nil
}
