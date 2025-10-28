package entity

import (
    "time"

    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    Username     string    `gorm:"type:text;not null;unique" json:"username" validate:"required,min=3,max=64"`
    PasswordHash string    `gorm:"type:text;not null" json:"password_hash" validate:"required,min=20"`
    CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (User) TableName() string { return "zenrows.user" }
