package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type DeviceProfile struct {
	ID            uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID        uuid.UUID         `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_name" json:"user_id" validate:"required"`
	TemplateID    *uuid.UUID        `gorm:"type:uuid" json:"template_id"`
	Name          string            `gorm:"type:text;not null;uniqueIndex:idx_user_name" json:"name" validate:"required,min=1,max=100"`
	DeviceType    string            `gorm:"type:text;not null" json:"device_type" validate:"required,oneof=desktop mobile"`
	Width         *int              `json:"width" validate:"omitempty,gt=0"`
	Height        *int              `json:"height" validate:"omitempty,gt=0"`
	UserAgent     *string           `gorm:"type:text" json:"user_agent" validate:"omitempty,min=1"`
	CountryCode   *string           `gorm:"type:char(2)" json:"country_code" validate:"omitempty,len=2,uppercase"`
	CustomHeaders datatypes.JSONMap `gorm:"type:jsonb" json:"custom_headers"`
	CreatedAt     time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

func (DeviceProfile) TableName() string { return "zenrows.device_profile" }
