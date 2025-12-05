package models

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	Phone        string    `json:"phone"`
	IDCard       string    `json:"id_card"`
	Department   string    `json:"department"`
	JobTitle     string    `json:"job_title"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Role       string `json:"role" binding:"omitempty"`
	Phone      string `json:"phone" binding:"omitempty"`
	IDCard     string `json:"id_card" binding:"omitempty"`
	Department string `json:"department" binding:"omitempty"`
	JobTitle   string `json:"job_title" binding:"omitempty"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Phone      string    `json:"phone"`
	IDCard     string    `json:"id_card"`
	Department string    `json:"department"`
	JobTitle   string    `json:"job_title"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
