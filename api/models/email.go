package models

type Email struct {
	YourName     string `json:"your_name" db:"your_name" binding:"required"`
	YourEmail    string `json:"from_email" db:"from_email" binding:"required,email"`
	YourPassword string `json:"your_api_key" db:"your_api_key" binding:"required,email"`
	ToEmail      string `json:"to_email" db:"to_email" binding:"required,email"`
	Title        string `json:"title" db:"title" binding:"required"`
	Text         string `json:"text" db:"text" binding:"required"`
}

type CreateEmail struct {
	YourName     string `json:"your_name" db:"your_name" binding:"required"`
	ToEmail      string `json:"to_email" db:"to_email" binding:"required,email"`
	Title        string `json:"title" db:"title" binding:"required"`
	Text         string `json:"text" db:"text" binding:"required"`
}
