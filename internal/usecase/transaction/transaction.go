package transaction

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"sync"
	"template/internal/model"
	"template/internal/repository"
	"template/internal/usecase"
)

type TransactionHandler struct {
	t repository.TransactionRepository
	u repository.UserRepository
	i repository.ItemRepository
}

func (t TransactionHandler) Purchase(ctx echo.Context, transaction model.RequestTransaction) error {
	var (
		err error
		mu  sync.Mutex
	)

	// handle potential race condition
	mu.Lock()
	defer mu.Unlock()

	tx, err := t.u.BeginTx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	item, err := t.i.GetItemByID(ctx, int64(transaction.ItemID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if item.ID == 0 {
		err = fmt.Errorf("not found items with id %v", transaction.ItemID)
		tx.Rollback()
		return err
	}

	if transaction.Quantity > item.Availability {
		err = fmt.Errorf("quantity is greater than avilability items")
		tx.Rollback()
		return err
	}

	if transaction.TotalPrice != int64(item.Price*transaction.Quantity) {
		err = fmt.Errorf("total price not matched, should equal with current price")
		tx.Rollback()
		return err
	}

	afterPurchase := model.RequestCreateItem{
		ID:           transaction.ItemID,
		Name:         item.Name,
		Rating:       item.Rating,
		Category:     item.Category,
		ImageURL:     item.ImageURL,
		Reputation:   item.Reputation,
		Price:        item.Price,
		Availability: item.Availability - transaction.Quantity,
		UserID:       int64(transaction.UserID),
		Value:        item.Value,
		IsActive:     true,
	}

	err = t.i.UpdateItem(ctx, afterPurchase, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	transaction.Availability = item.Availability

	err = t.t.Purchase(ctx, transaction, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func NewTransactionsUsecase(t repository.TransactionRepository, u repository.UserRepository, i repository.ItemRepository) usecase.TransactionUcase {
	return &TransactionHandler{t, u, i}
}
