package transaction

import (
	"template/internal/repository"
	"template/internal/usecase"
)

type TransactionHandler struct {
	t repository.TransactionRepository
	u repository.UserRepository
}

func NewTransactionsUsecase(t repository.TransactionRepository, u repository.UserRepository) usecase.TransactionUcase {
	return &TransactionHandler{t, u}
}
