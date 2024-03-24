package utils

import (
	"encoding/json"
	"log"
)

func StructToString(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshalling struct to JSON:", err)
		return ""
	}
	return string(jsonData)
}
