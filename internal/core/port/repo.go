package port

import "zenrows-challenge/internal/core/entity"

type UserRepo interface {
	RetrieveCredentials(u entity.User) (string, string, error)
}

type DeviceTemplateRepo interface {
	GetDeviceTemplates() ([]entity.DeviceTemplate, error)
}
