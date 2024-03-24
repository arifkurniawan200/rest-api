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
	Value        string `json:"value" db:"value"`
}

type RequestCreateItem struct {
	Name         string `json:"name"  validate:"required,blacklistWords,min=10"`
	Rating       int    `json:"rating"  validate:"required,inRange=0-5"`
	Category     string `json:"category"  validate:"required,checkCategory"`
	ImageURL     string `json:"image_url"  validate:"required,url"`
	Reputation   int    `json:"reputation"  validate:"required,inRange=0-1000"`
	Price        int    `json:"price"  validate:"required"`
	Availability int    `json:"availability"  validate:"required"`
	UserID       int64  `json:"-"`
	Value        string `json:"-"`
	ID           int    `json:"-"`
}

func (i *RequestCreateItem) GetReputationBadge() {
	switch {
	case i.Reputation <= 500:
		i.Value = "red"
	case i.Reputation <= 799:
		i.Value = "yellow"
	default:
		i.Value = "green"
	}
}

type Search struct {
	Rating       int    `json:"rating"`
	Category     string `json:"category"`
	Reputation   int    `json:"reputation"`
	Availability int    `json:"availability"`
}
