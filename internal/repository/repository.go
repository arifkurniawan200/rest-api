package repository

import (
	"database/sql"
	"template/internal/model"
)

type UserRepository interface {
	RegisterUser(user model.UserParam) error
	GetUserByEmail(email string) (model.User, error)
	BeginTx() (*sql.Tx, error)
}

type TransactionRepository interface {
}
