package transaction

import (
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	"template/internal/model"
	"template/internal/repository"
	"time"
)

type TransactionHandler struct {
	db *sql.DB
}

func (t TransactionHandler) Purchase(ctx echo.Context, transaction model.RequestTransaction, tx *sql.Tx) error {
	var (
		err error
	)

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	// check if using transaction
	if tx != nil {
		_, err = tx.ExecContext(ctxWithTimeout, insertPurchase, transaction.UserID, transaction.ItemID, transaction.TransactionType, transaction.TransactionDate, transaction.Quantity, transaction.TotalPrice, transaction.Availability)
		if err != nil {
			return err
		}
	} else {
		// using existing database from dependencies injection
		_, err = t.db.ExecContext(ctxWithTimeout, insertPurchase, transaction.UserID, transaction.ItemID, transaction.TransactionType, transaction.TransactionDate, transaction.Quantity, transaction.TotalPrice, transaction.Availability)
		if err != nil {
			return err
		}
	}
	return err
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &TransactionHandler{db}
}
