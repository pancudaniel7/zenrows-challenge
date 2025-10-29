package http

import (
	"zenrows-challenge/internal/core/entity"
)

func mapToDeviceTemplatesResponse(e entity.DeviceTemplate) DeviceTemplatesResponse {
	headers := make(map[string]string)
	for k, v := range e.DefaultHeaders {
		if str, ok := v.(string); ok {
			headers[k] = str
		}
	}

	return DeviceTemplatesResponse{
		ID:             e.ID,
		Name:           e.Name,
		DeviceType:     e.DeviceType,
		Width:          e.Width,
		Height:         e.Height,
		UserAgent:      e.UserAgent,
		CountryCode:    e.CountryCode,
		DefaultHeaders: headers,
		CreatedAt:      e.CreatedAt,
	}
}
