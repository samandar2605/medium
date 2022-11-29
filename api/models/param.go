package models

type GetAllParams struct {
	Limit      int    `json:"limit" binding:"required" default:"10"`
	Page       int    `json:"page" binding:"required" default:"1"`
	Search     string `json:"search"`
	PostId     int    `json:"post_id"`
	UserId     int    `json:"user_id"`
	SortByDate string `json:"sort_by_date" binding:"required,oneof=asc desc --" oneof:"" default:"desc"`
}

type GetAllResponse struct {
	Category Category
	Count    int
}

type VerifyUser struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}

type UserProfile struct {
	Id              int     `json:"id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	ProfileImageUrl *string `json:"profile_image_url"`
}
