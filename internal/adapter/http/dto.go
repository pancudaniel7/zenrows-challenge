package http

import (
	"time"

	"github.com/google/uuid"
)

type DeviceTemplatesResponse struct {
	ID             uuid.UUID         `json:"id"`
	Name           string            `json:"name"`
	DeviceType     string            `json:"device_type"`
	Width          *int              `json:"width,omitempty"`
	Height         *int              `json:"height,omitempty"`
	UserAgent      string            `json:"user_agent"`
	CountryCode    *string           `json:"country_code,omitempty"`
	DefaultHeaders map[string]string `json:"default_headers,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
}

type DeviceProfileResponse struct {
	ID            uuid.UUID         `json:"id"`
	UserID        uuid.UUID         `json:"user_id"`
	TemplateID    *uuid.UUID        `json:"template_id,omitempty"`
	Name          string            `json:"name"`
	DeviceType    string            `json:"device_type"`
	Width         *int              `json:"width,omitempty"`
	Height        *int              `json:"height,omitempty"`
	UserAgent     *string           `json:"user_agent,omitempty"`
	CountryCode   *string           `json:"country_code,omitempty"`
	CustomHeaders map[string]string `json:"custom_headers,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type DeviceProfileCreateRequest struct {
	TemplateID    *string           `json:"template_id,omitempty" validate:"omitempty,uuid4"`
	Name          string            `json:"name" validate:"required,min=1,max=100"`
	DeviceType    string            `json:"device_type" validate:"required,oneof=desktop mobile"`
	Width         *int              `json:"width,omitempty" validate:"omitempty,gt=0"`
	Height        *int              `json:"height,omitempty" validate:"omitempty,gt=0"`
	UserAgent     *string           `json:"user_agent,omitempty" validate:"omitempty,min=1"`
	CountryCode   *string           `json:"country_code,omitempty" validate:"omitempty,len=2,uppercase"`
	CustomHeaders map[string]string `json:"custom_headers,omitempty"`
}

type DeviceProfileUpdateRequest struct {
	TemplateID    *string           `json:"template_id,omitempty" validate:"omitempty,uuid4"`
	Name          *string           `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	DeviceType    *string           `json:"device_type,omitempty" validate:"omitempty,oneof=desktop mobile"`
	Width         *int              `json:"width,omitempty" validate:"omitempty,gt=0"`
	Height        *int              `json:"height,omitempty" validate:"omitempty,gt=0"`
	UserAgent     *string           `json:"user_agent,omitempty" validate:"omitempty,min=1"`
	CountryCode   *string           `json:"country_code,omitempty" validate:"omitempty,len=2,uppercase"`
	CustomHeaders map[string]string `json:"custom_headers,omitempty"`
}
