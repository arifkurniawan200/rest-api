package app

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the Authorization header from the request
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "Authorization header is missing",
				})
			}

			// Check if the Authorization header starts with "Bearer"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status":  401,
					"message": "Invalid authorization format",
				})
			}

			// Extract the token from the header
			tokenString := parts[1]

			// Parse and validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Invalid signing method")
				}
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status":  401,
					"message": "Invalid token",
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status":  401,
					"message": "Failed to parse claims",
				})
			}

			// Store the claims in the context for later use
			c.Set("auth", claims)

			return next(c)
		}
	}
}

// AdminMiddleware checks if the user is an admin based on JWT claims.
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the claims from the context
		claims, ok := c.Get("claims").(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "Invalid or missing claims",
			})
		}

		// Check if the "is_admin" claim exists and is true
		isAdmin, ok := claims["is_admin"].(bool)
		if !ok || !isAdmin {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "Access denied, only admin allowed to access this menu",
			})
		}

		// If the user is an admin, proceed to the next middleware or handler
		return next(c)
	}
}

// onlyJSONMiddleware check content type, only application json allowed
func onlyJSONMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get("Content-Type")
		if contentType != "application/json" && c.Request().Method == "POST" {
			return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"error": "Only application/json content type is allowed"})
		}
		return next(c)
	}
}
