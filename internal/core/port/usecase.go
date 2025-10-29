package port

import "zenrows-challenge/internal/core/entity"

type AuthenticationService interface {
	CheckCredentials(username string, password string) (string, error)
}

type DeviceTemplateService interface {
	RetrieveDeviceTemplates() ([]entity.DeviceTemplate, error)
}
