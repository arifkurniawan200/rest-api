package repository

import "database/sql"

type TransactionHandler struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &TransactionHandler{db}
}
