package repository

import "database/sql"

type historyHandler struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) HistoryRepository {
	return &historyHandler{db}
}
