package items

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"strings"
	"template/internal/model"
	"template/internal/repository"
	"time"
)

type ItemsHandler struct {
	db *sql.DB
}

func (i ItemsHandler) UpdateItem(ctx echo.Context, item model.RequestCreateItem, tx *sql.Tx) error {
	var (
		err error
	)

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	// check if using transaction
	if tx != nil {
		_, err = tx.ExecContext(ctxWithTimeout, updateItembyID, item.Name, item.Rating, item.Category, item.ImageURL, item.Reputation, item.Price, item.Availability, item.Value, item.IsActive, item.ID)
		if err != nil {
			return err
		}
	} else {
		// using existing database from dependencies injection
		_, err = i.db.ExecContext(ctxWithTimeout, updateItembyID, item.Name, item.Rating, item.Category, item.ImageURL, item.Reputation, item.Price, item.Availability, item.Value, item.IsActive, item.ID)
		if err != nil {
			return err
		}
	}
	return err
}

func (i ItemsHandler) CreateItem(ctx echo.Context, item model.RequestCreateItem) error {
	_, err := i.db.Exec(insertNewItem, item.Name, item.Rating, item.Category, item.ImageURL, item.Reputation, item.Price, item.Availability, item.Value)
	if err != nil {
		return err
	}
	return err
}

func (i ItemsHandler) GetItemByID(ctx echo.Context, itemID int64, tx *sql.Tx) (model.Item, error) {
	var (
		data model.Item
		rows *sql.Rows
		err  error
	)

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	// check if using transaction
	if tx != nil {
		rows, err = tx.QueryContext(ctxWithTimeout, geItemByID, itemID)
		if err != nil {
			return data, err
		}
	} else {
		// using existing database from dependencies injection
		rows, err = i.db.QueryContext(ctxWithTimeout, geItemByID, itemID)
		if err != nil {
			return data, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&data.ID, &data.Name, &data.Rating, &data.Category, &data.ImageURL, &data.Reputation, &data.Price, &data.Availability, &data.Value, &data.IsActive); err != nil {
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
		if err := rows.Scan(&t.ID, &t.Name, &t.Rating, &t.Category, &t.ImageURL, &t.Reputation, &t.Price, &t.Availability, &t.Value, &t.IsActive); err != nil {
			return nil, err
		}
		data = append(data, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func (i ItemsHandler) GetListPublicItem(ctx echo.Context, param model.Search) ([]model.Item, error) {
	var data []model.Item

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request().Context(), 5*time.Second)
	defer cancel()

	query := getListItem
	condition := []string{}
	values := []interface{}{}
	where := ""

	if param.Rating != 0 {
		condition = append(condition, "rating = ?")
		values = append(values, param.Rating)
	}
	if param.Category != "" {
		condition = append(condition, "category = ?")
		values = append(values, param.Category)
	}
	if param.Reputation != 0 {
		condition = append(condition, "reputation = ?")
		values = append(values, param.Reputation)
	}
	if param.Availability != 0 {
		condition = append(condition, "availability = ?")
		values = append(values, param.Rating)
	}

	if len(condition) > 0 {

		for i, x := range condition {
			condition[i] = strings.Replace(x, "?", fmt.Sprintf("$%v", i+1), -1)
		}
		where = strings.Join(condition, " AND ")
		query = query + "WHERE " + where
	}
	var rows *sql.Rows
	var err error

	if len(condition) > 0 {
		rows, err = i.db.QueryContext(ctxWithTimeout, query, values...)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = i.db.QueryContext(ctxWithTimeout, query)
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var t model.Item
		if err := rows.Scan(&t.ID, &t.Name, &t.Rating, &t.Category, &t.ImageURL, &t.Reputation, &t.Price, &t.Availability, &t.Value, &t.IsActive); err != nil {
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
