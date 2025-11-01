package http

import (
	"errors"
	"net/http"
	"zenrows-challenge/internal/pkg/apperr"

	"github.com/gofiber/fiber/v3"
)

// handleError maps application errors to HTTP responses for device templates.
func (h *DeviceTemplateHandlerImpl) handleError(c fiber.Ctx, err error) error {
	var (
		inv *apperr.InvalidArgErr
		nf  *apperr.NotFoundErr
		ae  *apperr.AlreadyExistsErr
		na  *apperr.NotAuthorizedErr
		in  *apperr.InternalErr
	)
	switch {
	case errors.As(err, &inv):
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"code": inv.Code(), "message": inv.Message()})
	case errors.As(err, &nf):
		return c.Status(http.StatusNotFound).JSON(map[string]string{"code": nf.Code(), "message": nf.Message()})
	case errors.As(err, &ae):
		return c.Status(http.StatusConflict).JSON(map[string]string{"code": ae.Code(), "message": ae.Message()})
	case errors.As(err, &na):
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"code": na.Code(), "message": na.Message()})
	case errors.As(err, &in):
		fallthrough
	default:
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"code": "INTERNAL_ERROR", "message": "internal error"})
	}
}
