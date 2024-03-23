package model

import "time"

// UserParam adalah model struct untuk request
type UserParam struct {
	NIK       string    `json:"NIK" validate:"required"`
	FullName  string    `json:"full_name" validate:"required"`
	Username  string    `json:"legal_name,omitempty" validate:"required"`
	BornPlace string    `json:"born_place,omitempty" validate:"required"`
	BornDate  time.Time `json:"born_date,omitempty" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required"`
}

// User adalah model struct untuk tabel user
type User struct {
	ID        int        `json:"id" db:"id"`
	NIK       string     `json:"NIK" db:"NIK"`
	FullName  string     `json:"full_name" db:"full_name"`
	BornPlace string     `json:"born_place,omitempty" db:"born_place"`
	BornDate  time.Time  `json:"born_date,omitempty" db:"born_date"`
	Email     string     `json:"email" db:"email"`
	IsAdmin   bool       `json:"is_admin,omitempty" db:"is_admin"`
	Password  string     `json:"password" db:"password"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
