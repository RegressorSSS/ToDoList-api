package model

import "time"

type Task struct {
	Id          int       `json: "id"`
	Title       string    `json: "title"`
	Description string    `json: "description"`
	IsComplete  bool      `json: "is_complete"`
	CreatedAt   time.Time `json: "created_at"`
}
