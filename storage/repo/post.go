package repo

import "time"

type GetPostQuery struct {
	Page       int    `json:"page" db:"page" binding:"required" default:"1"`
	Limit      int    `json:"limit" db:"limit" binding:"required" default:"10"`
	UserID     int    `json:"user_id"`
	CategoryID int    `json:"post_id"`
	SortByDate string `json:"sort_by_date" enums:"asc,desc" default:"desc"`
}

const (
	UserTypeSuperadmin = "superadmin"
	UserTypeUser       = "user"
)

type GetAllPostResult struct {
	Post  []*Post
	Count int
}

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

type UserProfile struct {
	Id              int     `json:"id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	ProfileImageUrl *string `json:"profile_image_url"`
}

type PostStorageI interface {
	Create(p *Post) (*Post, error)
	Get(id int) (*Post, error)
	GetAll(param GetPostQuery) (*GetAllPostResult, error)
	Update(usr *Post) (*Post, error)
	Delete(id int) error
	GetUserInfo(id int) (int)
	ViewsInc(id int) error
}
