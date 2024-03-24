package repository

import (
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	"template/internal/model"
	"time"
)

type ItemsHandler struct {
	db *sql.DB
}

func (i ItemsHandler) GetMyItem(ctx echo.Context, userID int64) ([]model.Item, error) {
	var data []model.Item

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	rows, err := i.db.QueryContext(ctxWithTimeout, getMyItems, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.Item
		if err := rows.Scan(&t.ID, &t.Name, &t.Rating, &t.Category, &t.ImageURL, &t.Reputation, &t.Price, &t.Availability); err != nil {
			return nil, err
		}
		data = append(data, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func (i ItemsHandler) GetListPublicItem(ctx echo.Context) ([]model.Item, error) {
	var data []model.Item

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	rows, err := i.db.QueryContext(ctxWithTimeout, getListItem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.Item
		if err := rows.Scan(&t.ID, &t.Name, &t.Rating, &t.Category, &t.ImageURL, &t.Reputation, &t.Price, &t.Availability); err != nil {
			return nil, err
		}
		data = append(data, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &ItemsHandler{db}
}
