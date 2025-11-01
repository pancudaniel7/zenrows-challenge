package port

import (
	"zenrows-challenge/internal/core/entity"

	"github.com/google/uuid"
)

// UserRepo exposes persistence operations for authenticating users.
type UserRepo interface {
	// RetrieveCredentials returns the existing password hash for the provided user.
	RetrieveCredentials(u entity.User) (string, string, error)
}

// DeviceTemplateRepo exposes queries for shared device templates.
type DeviceTemplateRepo interface {
	// GetDeviceTemplates returns every available device template.
	GetDeviceTemplates() ([]entity.DeviceTemplate, error)
	// GetDeviceTemplateByID retrieves a template by its identifier.
	GetDeviceTemplateByID(id *uuid.UUID) (*entity.DeviceTemplate, error)
}

// DeviceProfileRepo exposes CRUD operations for user device profiles.
type DeviceProfileRepo interface {
	// ListDeviceProfiles returns the paginated profiles for a given user.
	ListDeviceProfiles(userID string, page, pageSize int) ([]entity.DeviceProfile, error)
	// CreateDeviceProfile persists a new profile.
	CreateDeviceProfile(dp *entity.DeviceProfile) error
	// UpdateDeviceProfile modifies an existing profile.
	UpdateDeviceProfile(dp *entity.DeviceProfile) error
	// DeleteDeviceProfile removes a profile belonging to the supplied user.
	DeleteDeviceProfile(userID, id string) error
}
