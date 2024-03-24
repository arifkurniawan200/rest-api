package items

import (
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	"template/internal/model"
	"template/internal/repository"
	"time"
)

type ItemsHandler struct {
	db *sql.DB
}

func (i ItemsHandler) CreateItem(ctx echo.Context, item model.RequestCreateItem) error {
	_, err := i.db.Exec(insertNewItem, item.Name, item.Rating, item.Category, item.ImageURL, item.Reputation, item.Price, item.Availability, item.Value)
	if err != nil {
		return err
	}
	return err
}

func (i ItemsHandler) GetItemByID(ctx echo.Context, itemID int64) (model.Item, error) {
	var data model.Item

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	rows, err := i.db.QueryContext(ctxWithTimeout, geItemByID, itemID)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&data.ID, &data.Name, &data.Rating, &data.Category, &data.ImageURL, &data.Reputation, &data.Price, &data.Availability, &data.Value); err != nil {
			return data, err
		}
	}

	if err := rows.Err(); err != nil {
		return data, err
	}
	return data, nil
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
		if err := rows.Scan(&t.ID, &t.Name, &t.Rating, &t.Category, &t.ImageURL, &t.Reputation, &t.Price, &t.Availability, &t.Value); err != nil {
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
		if err := rows.Scan(&t.ID, &t.Name, &t.Rating, &t.Category, &t.ImageURL, &t.Reputation, &t.Price, &t.Availability, &t.Value); err != nil {
			return nil, err
		}
		data = append(data, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func NewItemRepository(db *sql.DB) repository.ItemRepository {
	return &ItemsHandler{db}
}
