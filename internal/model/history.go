package model

import "time"

type TableHistory struct {
	ID         int         `json:"id"`
	TableName  string      `json:"table_name"`
	TableKey   int         `json:"table_key"`
	DataBefore interface{} `json:"data_before"`
	DataAfter  interface{} `json:"data_after"`
	UserID     int64       `json:"user_id"`
	CreatedAt  time.Time   `json:"created_at"`
}
