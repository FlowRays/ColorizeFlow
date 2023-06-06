package model

import (
	"time"
)

type Image struct {
	ID        uint `gorm:"primaryKey"`
	Path      string
	UserID    uint
	Type      string
	CreatedAt time.Time
}

type ImageWithUser struct {
	ID        uint
	Username  string
	CreatedAt time.Time
}
