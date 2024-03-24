package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetSession(c echo.Context) (jwt.MapClaims, error) {
	auth, ok := c.Get("auth").(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed get mapcalims")
	}
	return auth, nil
}
