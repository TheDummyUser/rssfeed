package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string         `gorm:"size:50;unique;not null" json:"username"`
	Email        string         `gorm:"size:100;unique;not null" json:"email"`
	Password     string         `gorm:"not null" json:"-"`
	RefreshToken *string        `gorm:"type:text;unique" json:"refresh_token,omitempty"` // Nullable, unique, omitempty
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type SignupRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required"`
}

type RssData struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint   `gorm:"not null" json:"user_id"`
	Name   string `gorm:"not null;not null" json:"name"`
	Url    string `gorm:"not null;not null" json:"url"`
}

type RssRequest struct {
	Name string `json:"name" validate:"required"`
	Url  string `json:"url" validate:"required"`
}
