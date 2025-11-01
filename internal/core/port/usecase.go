package port

import (
	"context"
	"zenrows-challenge/internal/core/entity"
)

type AuthenticationService interface {
	CheckCredentials(username string, password string) (string, error)
}

type DeviceTemplateService interface {
	RetrieveDeviceTemplates() ([]entity.DeviceTemplate, error)
}

type DeviceProfileService interface {
	ListDeviceProfilesByUserID(ctx context.Context, page, pageSize int) ([]entity.DeviceProfile, error)
	CreateDeviceProfile(ctx context.Context, dp *entity.DeviceProfile) error
	UpdateDeviceProfile(ctx context.Context, dp *entity.DeviceProfile) (*entity.DeviceProfile, error)
	DeleteDeviceProfile(ctx context.Context, id string) error
}
