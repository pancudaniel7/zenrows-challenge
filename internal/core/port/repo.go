package port

import (
	"zenrows-challenge/internal/core/entity"

	"github.com/google/uuid"
)

type UserRepo interface {
	RetrieveCredentials(u entity.User) (string, string, error)
}

type DeviceTemplateRepo interface {
	GetDeviceTemplates() ([]entity.DeviceTemplate, error)
	GetDeviceTemplateByID(id *uuid.UUID) (*entity.DeviceTemplate, error)
}

type DeviceProfileRepo interface {
	ListDeviceProfiles(userID string, page, pageSize int) ([]entity.DeviceProfile, error)
	CreateDeviceProfile(dp *entity.DeviceProfile) error
	UpdateDeviceProfile(dp *entity.DeviceProfile) error
	DeleteDeviceProfile(userID, id string) error
}
