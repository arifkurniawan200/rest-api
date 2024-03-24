package items

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"template/internal/model"
	"template/internal/repository"
	"template/internal/utils"
	"time"
)

type historyHandler struct {
	db *sql.DB
}

func (h historyHandler) GetHistory(ctx echo.Context, itemID int64) ([]model.TableHistory, error) {
	var data []model.TableHistory

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	rows, err := h.db.QueryContext(ctxWithTimeout, getHistoryChanges, "items", itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.TableHistory
		var dataBefore string
		var dataAfter string
		if err := rows.Scan(&t.ID, &t.TableName, &t.TableKey, &dataBefore, &dataAfter, &t.UserID, &t.CreatedAt); err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(dataBefore), &t.DataBefore)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("Error unmarshalling JSON:", err)
		}
		err = json.Unmarshal([]byte(dataAfter), &t.DataAfter)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling JSON:", err)
		}

		data = append(data, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
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
