package item

import (
	"github.com/labstack/echo/v4"
	"template/internal/model"
	"template/internal/repository"
	"template/internal/usecase"
)

type ItemHandler struct {
	t repository.TransactionRepository
	u repository.UserRepository
	h repository.HistoryRepository
	i repository.ItemRepository
}

func (i ItemHandler) AddItem(ctx echo.Context, item model.RequestCreateItem) error {
	switch {
	case item.Reputation <= 500:
		item.Value = "red"
	case item.Reputation <= 799:
		item.Value = "yellow"
	default:
		item.Value = "green"
	}
	return i.i.CreateItem(ctx, item)
}

func (i ItemHandler) GetItemByItemID(ctx echo.Context, itemID int64) (model.Item, error) {
	return i.i.GetItemByID(ctx, itemID)
}

func (i ItemHandler) GetMyItem(ctx echo.Context, userID int64) ([]model.Item, error) {
	return i.i.GetMyItem(ctx, userID)
}

func (i ItemHandler) GetMarketItem(ctx echo.Context) ([]model.Item, error) {
	return i.i.GetListPublicItem(ctx)
}

func NewItemUsecase(t repository.TransactionRepository, u repository.UserRepository, h repository.HistoryRepository, i repository.ItemRepository) usecase.ItemUcase {
	return &ItemHandler{t, u, h, i}
}
