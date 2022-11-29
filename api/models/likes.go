package models

type Like struct {
	Id     int    `json:"id"`
	PostId int    `json:"post_id"`
	UserId int    `json:"user_id"`
	Status bool `json:"status"`
}

type CreateLike struct {
	PostId int    `json:"post_id"`
	Status bool `json:"status"`
}

type UpdateLike struct {
	Id     int    `json:"id"`
	PostId int    `json:"post_id"`
	Status bool `json:"status"`
}

type GetAllLikesParams struct {
	Limit  int `json:"limit" binding:"required" default:"10"`
	Page   int `json:"page" binding:"required" default:"1"`
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

type GetAllLikesResponse struct {
	Likes []*Like `json:"likes"`
	Count int     `json:"count"`
}
