package models

import (
	"time"
)

type User struct {
	Id              int       `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	PhoneNumber     *string    `json:"phone_number"`
	Email           string    `json:"email"`
	Gender          *string    `json:"gender"`
	Password        string    `json:"password"`
	Username        string    `json:"username"`
	ProfileImageUrl *string    `json:"profile_image_url"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateUser struct {
	FirstName       string `json:"first_name" binding:"required,min=2,max=30"`
	LastName        string `json:"last_name" binding:"required,min=2,max=30"`
	PhoneNumber     *string `json:"phone_number"`
	Email           string `json:"email" binding:"required,email"`
	Gender          *string `json:"gender" binding:"oneof=male female" default:"male"`
	Password        string `json:"password" binding:"required,min=6,max=16"`
	UserName        string `json:"username" binding:"required,min=2,max=30"`
	ProfileImageUrl *string `json:"profile_image_url"`
	Type            string `json:"type" binding:"required,oneof=superadmin user" default:"user"`
}

type GetAllUsersParams struct {
	Limit      int    `json:"limit" binding:"required" default:"10"`
	Page       int    `json:"page" binding:"required" default:"1"`
	Search     string `json:"search"`
	SortByDate string `json:"sort_by_date" enums:"asc,desc"`
}

type GetAllUsersResponse struct {
	Users []*User `json:"users"`
	Count int     `json:"count"`
}
