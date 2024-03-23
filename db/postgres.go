package db

import (
	"database/sql"
	"fmt"
	"sync"
	"template/config"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	err  error
	once sync.Once
)

func NewDatabase(cfg config.Database) (*sql.DB, error) {
	once.Do(func() {
		// Connection string for PostgreSQL
		dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name)

		// Open connection to PostgreSQL
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			return
		}

		if cfg.ActivePool {
			// Set the maximum and minimum connection pool size
			db.SetMaxIdleConns(cfg.MaxPool) // Maximum number of idle connections
			db.SetMaxOpenConns(cfg.MaxPool) // Maximum number of open connections
		}

		// Test the database connection
		if err = db.Ping(); err != nil {
			return
		}
	})

	return db, err
}
