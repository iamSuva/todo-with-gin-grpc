package models

import (
	"time"
)

type Task struct {
	Id            int       `json:"id"`
	Title         string    `json:"title"  validate:"required,min=5"`
	Description   string    `json:"description" validate:"required,min=5"`
	IsCompleted   bool      `json:"isCompleted"`
	CreatedAt_UTC time.Time `json:"createdAt_utc"`
	UpdatedAt_UTC time.Time `json:"updatedAt_utc"`
	UserId        int       `json:"userId"`
}

type User struct {
	UserId   int    `json:"userId"`
	Username string `json:"username" validate:"required,min=5"`
	Password string `json:"password" validate:"required,password"`
}
