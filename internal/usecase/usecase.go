package usecase

import (
	"github.com/labstack/echo/v4"
	"template/internal/model"
)

type UserUcase interface {
	RegisterCustomer(ctx echo.Context, customer model.UserParam) error
	GetUserInfoByEmail(ctx echo.Context, email string) (model.User, error)
}

type TransactionUcase interface {
}
