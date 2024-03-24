package app

import "github.com/golang-jwt/jwt/v5"

type jwtCustomClaims struct {
	Email string `json:"email"`
	ID    int64  `json:"id"`
	jwt.RegisteredClaims
}

type ResponseSuccess struct {
	Messages string      `json:"messages,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type ResponseFailed struct {
	Status   int         `json:"status"`
	Messages string      `json:"messages,omitempty"`
	Error    interface{} `json:"error,omitempty"`
}
