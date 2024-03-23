package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"template/config"
)

var (
	db  *sql.DB
	err error
)

func NewDatabase(cfg config.Database) (*sql.DB, error) {
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err = sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}

	if cfg.ActivePool {
		// Set the maximum and minimum connection pool size
		db.SetMaxIdleConns(cfg.MaxPool) // Maximum number of idle connections
		db.SetMaxOpenConns(cfg.MinPool) // Maximum number of open connections

	}
	// Test the database connection
	if err = db.Ping(); err != nil {
		return db, err
	}
	return db, err
}
