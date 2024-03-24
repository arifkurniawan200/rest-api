package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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
