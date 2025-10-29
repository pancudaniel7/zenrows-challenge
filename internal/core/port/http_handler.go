package port

import "github.com/gofiber/fiber/v3"

type DeviceTemplateHandler interface {
	List(c fiber.Ctx) error
}
