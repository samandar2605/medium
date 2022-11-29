package models

import "time"

type Post struct {
	Id          int         `json:"id" db:"id"`
	Title       string      `json:"title" db:"title"`
	Description string      `json:"description" db:"description"`
	ImageUrl    string      `json:"image_url" db:"image_url"`
	UserId      int         `json:"user_id" db:"user_id"`
	CategoryId  int         `json:"category_id" db:"category_id"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
	ViewsCount  int         `json:"views_count" db:"views_count"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	User        UserProfile `json:"user"`
}

type CreatePost struct {
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	ImageUrl    string `json:"image_url" db:"image_url"`
	CategoryId  int    `json:"category_id" db:"category_id"`
}

type GetAllPostsParams struct {
	Limit      int `json:"limit" binding:"required" default:"10"`
	Page       int `json:"page" binding:"required" default:"1"`
	UserID     int `json:"user_id"`
	CategoryId int `json:"category_id"`
}

type GetAllPostsResponse struct {
	Posts []*Post `json:"posts"`
	Count int     `json:"count"`
}
