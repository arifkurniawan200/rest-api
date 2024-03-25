package repository

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"template/internal/model"
)

type UserRepository interface {
	RegisterUser(user model.RequestRegisterUser) error
	GetUserByEmail(email string) (model.User, error)
	BeginTx() (*sql.Tx, error)
}

type TransactionRepository interface {
	Purchase(ctx echo.Context, transaction model.RequestTransaction, tx *sql.Tx) error
}

type ItemRepository interface {
	GetListPublicItem(ctx echo.Context, param model.Search) ([]model.Item, error)
	GetMyItem(ctx echo.Context, userID int64) ([]model.Item, error)
	GetItemByID(ctx echo.Context, itemID int64, tx *sql.Tx) (model.Item, error)
	CreateItem(ctx echo.Context, item model.RequestCreateItem) error
	UpdateItem(ctx echo.Context, item model.RequestCreateItem, tx *sql.Tx) error
}

type HistoryRepository interface {
	SaveHistory(ctx echo.Context, history model.TableHistory, tx *sql.Tx) error
	GetHistory(ctx echo.Context, itemID int64) ([]model.TableHistory, error)
}
