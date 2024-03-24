package items

import (
	"database/sql"
	"template/internal/repository"
)

type historyHandler struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) repository.HistoryRepository {
	return &historyHandler{db}
}
