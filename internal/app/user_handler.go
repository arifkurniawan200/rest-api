package app

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"template/internal/model"
	"template/internal/utils"
	"time"
)

func (u handler) RegisterUser(c echo.Context) error {
	var (
		user = new(model.RequestRegisterUser)
		err  error
	)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to register user",
			Error:    err.Error(),
		})
	}

	if err = cv.Validate(c, user); err != nil {
		errorResponse := ResponseFailed{
			Messages: "Validation Error",
			Status:   http.StatusBadRequest,
			Error:    _FormatValidationError(err),
		}
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	err = u.User.RegisterCustomer(c, *user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to register user",
			Error:    err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, ResponseSuccess{
		Messages: "success register user",
	})
}

func (u handler) LoginUser(c echo.Context) error {
	var (
		user = new(model.RequestLogin)
		err  error
	)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to register user",
			Error:    err.Error(),
		})
	}

	if err = cv.Validate(c, user); err != nil {
		errorResponse := ResponseFailed{
			Messages: "Validation Error",
			Status:   http.StatusBadRequest,
			Error:    _FormatValidationError(err),
		}
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	userInfo, err := u.User.GetUserInfoByEmail(c, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "failed to login",
			Error:    err.Error(),
		})
	}

	if !utils.VerifyPassword(user.Password, userInfo.Password) {
		return c.JSON(http.StatusUnauthorized, ResponseFailed{
			Status:   http.StatusUnauthorized,
			Messages: "invalid username/password",
			Error:    "username or password is mismatch",
		})
	}
	claims := &jwtCustomClaims{
		userInfo.Email,
		int64(userInfo.ID),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	accessToken, err := token.SignedString([]byte(u.Cfg.Env.SecretKey))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{
			Status:   http.StatusInternalServerError,
			Messages: "error when generate token",
			Error:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, ResponseSuccess{
		Messages: "success login",
		Data: map[string]string{
			"token": accessToken,
		},
	})
}
