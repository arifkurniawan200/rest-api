package model

import "time"

type RequestTransaction struct {
	UserID          int       `json:"-"`
	ItemID          int       `json:"item_id" validate:"required"`
	TransactionType string    `json:"transaction_type" validate:"required"`
	TransactionDate time.Time `json:"-"`
	Quantity        int       `json:"quantity" validate:"required"`
	TotalPrice      int64     `json:"total_price" validate:"required"`
	Availability    int       `json:"-"`
}
