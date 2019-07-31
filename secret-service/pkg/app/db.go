package app

import (
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
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
	var err error
	for i := 1; i <= 10; i++ {
		db, err := gorm.Open("postgres", cfg.GetPostgresConnectionString())
		if err != nil {
			log.Warnf("Database not ready, retry after 1 second")
			time.Sleep(time.Second)
		} else {
			return db, err
		}
	}
	return nil, err
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
