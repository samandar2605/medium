package models

import "time"

type Category struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCategory struct {
	Title string `json:"title" binding:"required"`
}


type GetAllCategoriesResponse struct {
	Categories []*Category `json:"categories"`
	Count      int         `json:"count"`
}

type GetAllCategoryParams struct {
	Limit  int    `json:"limit" binding:"required" default:"10"`
	Page   int    `json:"page" binding:"required" default:"1"`
	Search string `json:"search"`
}
