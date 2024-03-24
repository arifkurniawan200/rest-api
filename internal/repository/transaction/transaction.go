package transaction

import (
	"database/sql"
	"template/internal/repository"
)

type TransactionHandler struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &TransactionHandler{db}
}
