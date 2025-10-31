package port

import "zenrows-challenge/internal/core/entity"

type UserRepo interface {
	RetrieveCredentials(u entity.User) (string, string, error)
}

type DeviceTemplateRepo interface {
    GetDeviceTemplates() ([]entity.DeviceTemplate, error)
}

type DeviceProfileRepo interface {
    ListDeviceProfiles(userID string, page, pageSize int) ([]entity.DeviceProfile, error)
    CreateDeviceProfile(dp *entity.DeviceProfile) error
    UpdateDeviceProfileSelective(dp *entity.DeviceProfile) error
    GetDeviceProfileByID(id string) (*entity.DeviceProfile, error)
    UpdateDeviceProfile(dp *entity.DeviceProfile) error
    DeleteDeviceProfile(id string) error
}
