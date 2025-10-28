package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type DeviceTemplate struct {
	ID             uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name           string            `gorm:"type:text;not null" json:"name" validate:"required,min=1,max=100"`
	DeviceType     string            `gorm:"type:text;not null" json:"device_type" validate:"required,oneof=desktop mobile"`
	Width          *int              `json:"width" validate:"omitempty,gt=0"`
	Height         *int              `json:"height" validate:"omitempty,gt=0"`
	UserAgent      string            `gorm:"type:text;not null" json:"user_agent" validate:"required,min=1"`
	CountryCode    *string           `gorm:"type:char(2)" json:"country_code" validate:"omitempty,len=2,uppercase"`
	DefaultHeaders datatypes.JSONMap `gorm:"type:jsonb" json:"default_headers"`
	CreatedAt      time.Time         `gorm:"autoCreateTime" json:"created_at"`
}

func (DeviceTemplate) TableName() string { return "zenrows.device_template" }
