package items

import (
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	"template/internal/model"
	"template/internal/repository"
	"template/internal/utils"
	"time"
)

type historyHandler struct {
	db *sql.DB
}

func (h historyHandler) SaveHistory(ctx echo.Context, history model.TableHistory, tx *sql.Tx) error {
	var (
		err error
	)

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	// check if using transaction
	if tx != nil {
		_, err = tx.ExecContext(ctxWithTimeout, saveHistoryChanges, history.TableName, history.TableKey, utils.StructToString(history.DataBefore), utils.StructToString(history.DataAfter), history.UserID, time.Now())
		if err != nil {
			return err
		}
	} else {
		// using existing database from dependencies injection
		_, err = h.db.ExecContext(ctxWithTimeout, saveHistoryChanges, history.TableName, history.TableKey, utils.StructToString(history.DataBefore), utils.StructToString(history.DataAfter), history.UserID, time.Now())
		if err != nil {
			return err
		}
	}
	return err
}

func NewHistoryRepository(db *sql.DB) repository.HistoryRepository {
	return &historyHandler{db}
}
