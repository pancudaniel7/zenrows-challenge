package port

import "github.com/gofiber/fiber/v3"

type DeviceTemplateHandler interface {
	List(c fiber.Ctx) error
}

type DeviceProfileHandler interface {
	ListDeviceProfilesByUserID(c fiber.Ctx) error
	CreateDeviceProfile(c fiber.Ctx) error
	GetDeviceProfileByID(c fiber.Ctx) error
	UpdateDeviceProfile(c fiber.Ctx) error
	DeleteDeviceProfile(c fiber.Ctx) error
}
