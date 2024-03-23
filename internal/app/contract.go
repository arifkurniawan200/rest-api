package app

import "github.com/golang-jwt/jwt/v5"

type jwtCustomClaims struct {
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	ID      int64  `json:"id"`
	jwt.RegisteredClaims
}

type ResponseSuccess struct {
	Messages string      `json:"messages,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type ResponseFailed struct {
	Messages string      `json:"messages,omitempty"`
	Error    interface{} `json:"error,omitempty"`
}
