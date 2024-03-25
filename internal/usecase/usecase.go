package usecase

import (
	"github.com/labstack/echo/v4"
	"template/internal/model"
)

type UserUcase interface {
	RegisterCustomer(ctx echo.Context, customer model.RequestRegisterUser) error
	GetUserInfoByEmail(ctx echo.Context, email string) (model.User, error)
}

type TransactionUcase interface {
	Purchase(ctx echo.Context, transaction model.RequestTransaction) error
}

type ItemUcase interface {
	GetMarketItem(ctx echo.Context, param model.Search) ([]model.Item, error)
	GetMyItem(ctx echo.Context, userID int64) ([]model.Item, error)
	GetItemByItemID(ctx echo.Context, itemID int64) (model.Item, error)
	AddItem(ctx echo.Context, item model.RequestCreateItem) error
	UpdateItem(ctx echo.Context, item model.RequestCreateItem) error
	DeleteItem(ctx echo.Context, userID, itemID int64) error
}
