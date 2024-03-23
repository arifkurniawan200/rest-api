package usecase

import "template/internal/repository"

type TransactionHandler struct {
	t repository.TransactionRepository
	u repository.UserRepository
}

func NewTransactionsUsecase(t repository.TransactionRepository, u repository.UserRepository) TransactionUcase {
	return &TransactionHandler{t, u}
}
