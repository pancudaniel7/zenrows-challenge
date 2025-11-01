package port

import "github.com/gofiber/fiber/v3"

// DeviceTemplateHandler defines the HTTP handlers for device template operations.
type DeviceTemplateHandler interface {
	// List returns the available device templates.
	List(c fiber.Ctx) error
}

// DeviceProfileHandler defines the HTTP handlers for device profile operations.
type DeviceProfileHandler interface {
	// ListDeviceProfilesByUserID returns profiles owned by the authenticated user.
	ListDeviceProfilesByUserID(c fiber.Ctx) error
	// CreateDeviceProfile persists a new device profile.
	CreateDeviceProfile(c fiber.Ctx) error
	// UpdateDeviceProfile modifies an existing device profile.
	UpdateDeviceProfile(c fiber.Ctx) error
	// DeleteDeviceProfile removes a device profile.
	DeleteDeviceProfile(c fiber.Ctx) error
}
