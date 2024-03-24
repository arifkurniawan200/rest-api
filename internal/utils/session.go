package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetSession(c echo.Context) (jwt.MapClaims, error) {
	auth, ok := c.Get("userAuth").(jwt.MapClaims)
	if !ok || auth == nil {
		return nil, fmt.Errorf("failed to get jwt.MapClaims from context")
	}
	return auth, nil
}
