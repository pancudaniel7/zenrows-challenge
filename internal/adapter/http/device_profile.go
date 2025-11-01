package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"zenrows-challenge/internal/core/port"
	"zenrows-challenge/internal/pkg/applog"
	"zenrows-challenge/internal/pkg/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeviceProfileHandlerImpl struct {
	log applog.AppLogger
	svc port.DeviceProfileService
	v   *validator.Validate
}

func NewDeviceProfileHandlerImpl(log applog.AppLogger, svc port.DeviceProfileService, v *validator.Validate) *DeviceProfileHandlerImpl {
	return &DeviceProfileHandlerImpl{log: log, svc: svc, v: v}
}

func (h *DeviceProfileHandlerImpl) ListDeviceProfilesByUserID(c fiber.Ctx) error {
	page, pageSize, err := parsePagination(c.Query("page", "1"), c.Query("page_size", "20"))
	if err != nil {
		return badRequest(c, err.Error())
	}

	ctx, _, err := h.userContext(c)
	if err != nil {
		return err
	}

	items, err := h.svc.ListDeviceProfilesByUserID(ctx, page, pageSize)
	if err != nil {
		return handleError(c, err)
	}

	resp := make([]DeviceProfileResponse, len(items))
	for i, item := range items {
		resp[i] = mapToDeviceProfileResponse(item)
	}
	return c.JSON(resp)
}

func (h *DeviceProfileHandlerImpl) CreateDeviceProfile(c fiber.Ctx) error {
	var req DeviceProfileCreateRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if err := h.v.Struct(req); err != nil {
		return badRequest(c, fmt.Sprintf("validation failed: %v", err))
	}

	ctx, userIDStr, err := h.userContext(c)
	if err != nil {
		return err
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return badRequest(c, "invalid user id in context")
	}

	dp, err := mapDeviceProfileCreateRequestToEntity(req, userUUID)
	if err != nil {
		return badRequest(c, err.Error())
	}

	if err := h.svc.CreateDeviceProfile(ctx, dp); err != nil {
		return handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(mapToDeviceProfileResponse(*dp))
}

func (h *DeviceProfileHandlerImpl) UpdateDeviceProfile(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return badRequest(c, "invalid device profile id")
	}

	var req DeviceProfileUpdateRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if isEmptyUpdateRequest(req) {
		return badRequest(c, "no fields supplied for update")
	}

	if err := h.v.Struct(req); err != nil {
		return badRequest(c, fmt.Sprintf("validation failed: %v", err))
	}

	ctx, userIDStr, err := h.userContext(c)
	if err != nil {
		return err
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return badRequest(c, "invalid user id in context")
	}

	dp := mapDeviceProfileUpdateRequestToEntity(req, id, userUUID)
	updated, err := h.svc.UpdateDeviceProfile(ctx, &dp)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(mapToDeviceProfileResponse(*updated))
}

func (h *DeviceProfileHandlerImpl) DeleteDeviceProfile(c fiber.Ctx) error {
	idStr := c.Params("id")
	if _, err := uuid.Parse(idStr); err != nil {
		return badRequest(c, "invalid device profile id")
	}

	ctx, _, err := h.userContext(c)
	if err != nil {
		return err
	}

	if err := h.svc.DeleteDeviceProfile(ctx, idStr); err != nil {
		return handleError(c, err)
	}

	return c.SendStatus(http.StatusNoContent)
}

func (h *DeviceProfileHandlerImpl) userContext(c fiber.Ctx) (context.Context, string, error) {
	val := c.Locals(middleware.AuthUserIDKey)
	userID, ok := val.(string)
	if !ok || userID == "" {
		return nil, "", c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"code":    "NOT_AUTHORIZED",
			"message": "unauthorized",
		})
	}
	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, userID)
	return ctx, userID, nil
}

func parsePagination(pageStr, sizeStr string) (int, int, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, 0, errors.New("invalid page parameter")
	}
	pageSize, err := strconv.Atoi(sizeStr)
	if err != nil || pageSize < 1 {
		return 0, 0, errors.New("invalid page_size parameter")
	}
	return page, pageSize, nil
}

func badRequest(c fiber.Ctx, msg string) error {
	return c.Status(http.StatusBadRequest).JSON(map[string]string{
		"code":    "INVALID_ARGUMENT",
		"message": msg,
	})
}

func isEmptyUpdateRequest(req DeviceProfileUpdateRequest) bool {
	return req.TemplateID == nil &&
		req.Name == nil &&
		req.DeviceType == nil &&
		req.Width == nil &&
		req.Height == nil &&
		req.UserAgent == nil &&
		req.CountryCode == nil &&
		req.CustomHeaders == nil
}
