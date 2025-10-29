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
