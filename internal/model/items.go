package model

type Item struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Rating       int    `json:"rating" db:"rating"`
	Category     string `json:"category" db:"category"`
	ImageURL     string `json:"image_url" db:"image_url"`
	Reputation   int    `json:"reputation" db:"reputation"`
	Price        int    `json:"price" db:"price"`
	Availability int    `json:"availability" db:"availability"`
}
