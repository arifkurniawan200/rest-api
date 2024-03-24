package model

import "time"

type BaseUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserParam adalah model struct untuk request
type RequestRegisterUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"omitempty"`
	BaseUser
}

type RequestLogin = BaseUser

// User adalah model struct untuk tabel user
type User struct {
	ID        int        `json:"id" db:"id"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"password" db:"password"`
	Type      string     `json:"type" db:"type"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
