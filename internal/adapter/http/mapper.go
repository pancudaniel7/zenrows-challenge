package http

import (
	"zenrows-challenge/internal/core/entity"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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

func mapToDeviceProfileResponse(e entity.DeviceProfile) DeviceProfileResponse {
	headers := make(map[string]string)
	for k, v := range e.CustomHeaders {
		if str, ok := v.(string); ok {
			headers[k] = str
		}
	}
	return DeviceProfileResponse{
		ID:            e.ID,
		UserID:        e.UserID,
		TemplateID:    e.TemplateID,
		Name:          e.Name,
		DeviceType:    e.DeviceType,
		Width:         e.Width,
		Height:        e.Height,
		UserAgent:     e.UserAgent,
		CountryCode:   e.CountryCode,
		CustomHeaders: headers,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

// mapDeviceProfileUpdateRequestToEntity converts a partial update DTO into a
// DeviceProfile entity suitable for selective updates. Only provided fields are
// set on the returned entity.
func mapDeviceProfileUpdateRequestToEntity(req DeviceProfileUpdateRequest, id uuid.UUID, userID uuid.UUID) entity.DeviceProfile {
	dp := entity.DeviceProfile{ID: id, UserID: userID}
	if req.TemplateID != nil {
		if *req.TemplateID == "" {
			dp.TemplateID = nil
		} else {
			tid, _ := uuid.Parse(*req.TemplateID)
			dp.TemplateID = &tid
		}
	}
	if req.Name != nil {
		dp.Name = *req.Name
	}
	if req.DeviceType != nil {
		dp.DeviceType = *req.DeviceType
	}
	if req.Width != nil {
		dp.Width = req.Width
	}
	if req.Height != nil {
		dp.Height = req.Height
	}
	if req.UserAgent != nil {
		dp.UserAgent = req.UserAgent
	}
	if req.CountryCode != nil {
		dp.CountryCode = req.CountryCode
	}
	if req.CustomHeaders != nil {
		headers := datatypes.JSONMap{}
		for k, v := range req.CustomHeaders {
			headers[k] = v
		}
		dp.CustomHeaders = headers
	}
	return dp
}
