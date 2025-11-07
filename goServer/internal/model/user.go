package model

import (
	"time"
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `json:"-"` // omit in responses
	Role      Role      `gorm:"type:VARCHAR(20);not null;default:'CLIENT'" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
