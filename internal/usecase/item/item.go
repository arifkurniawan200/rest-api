package item

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"sync"
	"template/internal/model"
	"template/internal/repository"
	"template/internal/usecase"
	"time"
)

type ItemHandler struct {
	t repository.TransactionRepository
	u repository.UserRepository
	h repository.HistoryRepository
	i repository.ItemRepository
}

func (i ItemHandler) UpdateItem(ctx echo.Context, item model.RequestCreateItem) error {
	var (
		err error
	)

	item.GetReputationBadge()

	tx, err := i.u.BeginTx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	fmt.Println("PASS")
	dataBefore, err := i.i.GetItemByID(ctx, int64(item.ID), tx)
	if err != nil {
		return err
	}

	err = i.i.UpdateItem(ctx, item, tx)
	if err != nil {
		return err
	}

	err = i.h.SaveHistory(ctx, model.TableHistory{
		TableName:  "items",
		TableKey:   item.ID,
		DataBefore: dataBefore,
		DataAfter:  item,
		CreatedAt:  time.Now(),
		UserID:     item.UserID,
	}, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (i ItemHandler) AddItem(ctx echo.Context, item model.RequestCreateItem) error {
	item.GetReputationBadge()
	return i.i.CreateItem(ctx, item)
}

func (i ItemHandler) GetItemByItemID(ctx echo.Context, itemID int64) (model.Item, error) {
	var (
		wg        sync.WaitGroup
		item      model.Item
		err       error
		histories []model.TableHistory
	)

	wg.Add(2)

	go func() error {
		defer wg.Done()
		item, err = i.i.GetItemByID(ctx, itemID, nil)
		if err != nil {
			return nil
		}
		return nil
	}()

	go func() error {
		defer wg.Done()
		histories, err = i.h.GetHistory(ctx, itemID)
		if err != nil {
			return err
		}
		return nil
	}()

	wg.Wait()

	item.HistoryChanges = histories

	return item, nil
}

func (i ItemHandler) GetMyItem(ctx echo.Context, userID int64) ([]model.Item, error) {
	return i.i.GetMyItem(ctx, userID)
}

func (i ItemHandler) GetMarketItem(ctx echo.Context, param model.Search) ([]model.Item, error) {
	return i.i.GetListPublicItem(ctx, param)
}

func NewItemUsecase(t repository.TransactionRepository, u repository.UserRepository, h repository.HistoryRepository, i repository.ItemRepository) usecase.ItemUcase {
	return &ItemHandler{t, u, h, i}
}
