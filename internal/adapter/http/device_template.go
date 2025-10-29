package http

import (
	"zenrows-challenge/internal/core/port"
	"zenrows-challenge/internal/pkg/applog"

	"github.com/gofiber/fiber/v3"
)

type DeviceTemplateHandlerImpl struct {
	log applog.AppLogger
	svc port.DeviceTemplateService
}

func NewDeviceTemplateHandlerImpl(log applog.AppLogger, svc port.DeviceTemplateService) *DeviceTemplateHandlerImpl {
	return &DeviceTemplateHandlerImpl{log: log, svc: svc}
}

func (h *DeviceTemplateHandlerImpl) List(c fiber.Ctx) error {
	h.log.Trace("DeviceTemplateHandlerImpl: List called")
	items, err := h.svc.RetrieveDeviceTemplates()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(items)
}
