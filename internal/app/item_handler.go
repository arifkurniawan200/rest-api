package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"template/internal/model"
	"template/internal/utils"
	"time"
)

func (u handler) ListItems(c echo.Context) error {
	var (
		item = new(model.Search)
		err  error
	)

	if err := c.Bind(item); err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to get item list",
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

	data, err := u.Items.GetMarketItem(c, *item)
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
			Messages: "failed to add item",
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

func (u handler) UpdateItem(c echo.Context) error {
	var (
		item = new(model.RequestCreateItem)
		err  error
	)

	reqID := c.QueryParam("id")
	if reqID == "" {
		return c.JSON(http.StatusBadRequest, ResponseFailed{
			Status:   http.StatusBadRequest,
			Messages: "failed to update item",
			Error:    "param id is required",
		})
	}

	if err := c.Bind(item); err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to update item",
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

	userId, ok := auth["id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "internal server error",
			Error:    "failed to get user id",
		})
	}

	id, _ := strconv.Atoi(reqID)
	item.UserID = int64(userId)
	item.ID = id

	err = u.Items.UpdateItem(c, *item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to update item",
			Error:    err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, ResponseSuccess{
		Messages: "success update item",
	})
}

func (u handler) DeleteItem(c echo.Context) error {
	var (
		err error
	)

	reqID := c.QueryParam("id")
	if reqID == "" {
		return c.JSON(http.StatusBadRequest, ResponseFailed{
			Status:   http.StatusBadRequest,
			Messages: "failed to delete item",
			Error:    "param id is required",
		})
	}

	auth, err := utils.GetSession(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseFailed{
			Status:   http.StatusUnauthorized,
			Messages: "invalid token",
			Error:    "access token is invalid or expired",
		})
	}

	userId, ok := auth["id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "internal server error",
			Error:    "failed to get user id",
		})
	}

	itemID, _ := strconv.Atoi(reqID)

	err = u.Items.DeleteItem(c, int64(userId), int64(itemID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to delete item",
			Error:    err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, ResponseSuccess{
		Messages: "success delete item",
	})
}

func (u handler) BuyItem(c echo.Context) error {
	var (
		transaction = new(model.RequestTransaction)
		err         error
	)

	if err := c.Bind(transaction); err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to buy item",
			Error:    err.Error(),
		})
	}

	if err = cv.Validate(c, transaction); err != nil {
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
	transaction.UserID = int(id)
	transaction.TransactionDate = time.Now()

	err = u.Transaction.Purchase(c, *transaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to buy item",
			Error:    err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, ResponseSuccess{
		Messages: "success buy item",
	})
}
