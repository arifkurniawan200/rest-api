package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"template/internal/model"
	"template/internal/utils"
)

func (u handler) ListItems(c echo.Context) error {
	data, err := u.Items.GetMarketItem(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to get item list",
			Error:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, ResponseSuccess{
		Messages: "success fetch item list",
		Data:     data,
	})
}

func (u handler) ListMyItems(c echo.Context) error {
	auth, err := utils.GetSession(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseFailed{
			Status:   http.StatusUnauthorized,
			Messages: "invalid token",
			Error:    "access token is invalid or expired",
		})
	}

	id, ok := auth["id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "internal server error",
			Error:    "failed to get user id",
		})
	}

	data, err := u.Items.GetMyItem(c, int64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to get my item list",
			Error:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, ResponseSuccess{
		Messages: "success fetch my item list",
		Data:     data,
	})
}

func (u handler) GetItemsByID(c echo.Context) error {
	id := c.QueryParam("id")
	itemID, _ := strconv.Atoi(id)
	data, err := u.Items.GetItemByItemID(c, int64(itemID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to get item",
			Error:    err.Error(),
		})
	}

	if data.ID == 0 {
		return c.JSON(http.StatusNotFound, ResponseFailed{
			Status:   http.StatusNotFound,
			Messages: "not found item",
			Error:    fmt.Errorf("not found item with id %v", itemID),
		})
	}
	return c.JSON(http.StatusOK, ResponseSuccess{
		Messages: "success fetch item",
		Data:     data,
	})
}

func (u handler) AddItem(c echo.Context) error {
	var (
		item = new(model.RequestCreateItem)
		err  error
	)

	if err := c.Bind(item); err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to register user",
			Error:    err.Error(),
		})
	}

	if err = cv.Validate(c, item); err != nil {
		errorResponse := ResponseFailed{
			Messages: "Validation Error",
			Status:   http.StatusBadRequest,
			Error:    _FormatValidationError(err),
		}
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	auth, err := utils.GetSession(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseFailed{
			Status:   http.StatusUnauthorized,
			Messages: "invalid token",
			Error:    "access token is invalid or expired",
		})
	}

	id, ok := auth["id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "internal server error",
			Error:    "failed to get user id",
		})
	}
	item.UserID = int64(id)

	err = u.Items.AddItem(c, *item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to add item",
			Error:    err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, ResponseSuccess{
		Messages: "success add item",
	})
}
