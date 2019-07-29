package app

import (
	"github.com/jinzhu/gorm"
)

// TODO move it to sql files
var createTableStatements = []string{
	`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,

	`CREATE TABLE IF NOT EXISTS secrets (
		hash uuid primary key DEFAULT uuid_generate_v4() NOT NULL,
		secret_text text,
		created_at timestamp with time zone,
		expires_at timestamp with time zone,
		remaining_views integer
	)`,
}

type DBConfig interface {
	GetPostgresConnectionString() string
}

func OpenDB(cfg DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", cfg.GetPostgresConnectionString())
	return db, err
}

func SetupDB(db *gorm.DB) error {
	return createTables(db)
}

func createTables(db *gorm.DB) error {
	for _, sql := range createTableStatements {
		_, err := db.DB().Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}
