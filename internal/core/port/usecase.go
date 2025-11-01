package port

import (
	"context"
	"zenrows-challenge/internal/core/entity"
)

// AuthenticationService exposes the business logic for verifying user credentials.
type AuthenticationService interface {
	// CheckCredentials returns the user ID when the supplied username and password are valid.
	CheckCredentials(username string, password string) (string, error)
}

// DeviceTemplateService exposes the use cases for shared device templates.
type DeviceTemplateService interface {
	// RetrieveDeviceTemplates returns every available device template.
	RetrieveDeviceTemplates() ([]entity.DeviceTemplate, error)
}

// DeviceProfileService exposes the use cases for user device profiles.
type DeviceProfileService interface {
	// ListDeviceProfilesByUserID returns paginated profiles scoped to the authenticated user.
	ListDeviceProfilesByUserID(ctx context.Context, page, pageSize int) ([]entity.DeviceProfile, error)
	// CreateDeviceProfile persists a new profile instance.
	CreateDeviceProfile(ctx context.Context, dp *entity.DeviceProfile) error
	// UpdateDeviceProfile applies modifications to an existing profile.
	UpdateDeviceProfile(ctx context.Context, dp *entity.DeviceProfile) (*entity.DeviceProfile, error)
	// DeleteDeviceProfile removes a profile by identifier.
	DeleteDeviceProfile(ctx context.Context, id string) error
}
