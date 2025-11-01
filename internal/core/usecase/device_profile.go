package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/core/port"
	"zenrows-challenge/internal/pkg/apperr"
	"zenrows-challenge/internal/pkg/applog"
	"zenrows-challenge/internal/pkg/middleware"
)

// DeviceProfileServiceImpl provides application logic for device profiles.
type DeviceProfileServiceImpl struct {
	log                applog.AppLogger
	repo               port.DeviceProfileRepo
	deviceTemplateRepo port.DeviceTemplateRepo
	v                  *validator.Validate
}

// NewDeviceProfileServiceImpl constructs a new DeviceProfileServiceImpl with the provided logger and repository.
func NewDeviceProfileServiceImpl(log applog.AppLogger, r port.DeviceProfileRepo, dtr port.DeviceTemplateRepo, v *validator.Validate) *DeviceProfileServiceImpl {
	return &DeviceProfileServiceImpl{log: log, repo: r, deviceTemplateRepo: dtr, v: v}
}

func (s *DeviceProfileServiceImpl) CreateDeviceProfile(ctx context.Context, dp *entity.DeviceProfile) error {
	s.log.Trace("device_profile.create", "user_id", dp.UserID.String(), "name", dp.Name)
	if err := s.v.Struct(dp); err != nil {
		return apperr.NewInvalidArgErr("invalid payload", err)
	}

	userID := ctx.Value(middleware.AuthUserIDKey).(string)
	if dp.TemplateID != nil {
		t, err := s.deviceTemplateRepo.GetDeviceTemplateByID(dp.TemplateID)
		if err != nil {
			return apperr.NewNotFoundErr("device template not found", err)
		}

		dp = s.createDeviceProfileByTemplate(t)
		uid, err := uuid.Parse(userID)
		if err != nil {
			return apperr.NewInvalidArgErr("invalid user id", err)
		}
		dp.UserID = uid
	}

	if err := s.repo.CreateDeviceProfile(dp); err != nil {
		s.log.Error("device_profile.create failed: %v", err)
		return mapRepoErr("create device profile", err)
	}
	return nil
}

func (s *DeviceProfileServiceImpl) ListDeviceProfilesByUserID(ctx context.Context, page, pageSize int) ([]entity.DeviceProfile, error) {
	s.log.Trace("device_profile.list_by_user_id", "page", page, "page_size", pageSize)

	userID := ctx.Value(middleware.AuthUserIDKey).(string)

	items, err := s.repo.ListDeviceProfiles(userID, page, pageSize)
	if err != nil {
		s.log.Error("device_profile.list failed: %v", err)
		return nil, mapRepoErr("list device profiles", err)
	}
	return items, nil
}

func (s *DeviceProfileServiceImpl) UpdateDeviceProfile(ctx context.Context, dp *entity.DeviceProfile) (*entity.DeviceProfile, error) {
	s.log.Trace("device_profile.update", "id", dp.ID.String(), "name", dp.Name)

	userID := dp.UserID.String()
	uid := ctx.Value(middleware.AuthUserIDKey).(string)
	if userID == "" || uid != userID {
		return nil, apperr.NewNotAuthorizedErr("NotAuthorizedErr", nil)
	}

	if err := s.v.Var(dp.ID, "required,uuid4"); err != nil {
		return nil, apperr.NewInvalidArgErr("invalid id", err)
	}
	if err := s.v.Var(dp.UserID, "required,uuid4"); err != nil {
		return nil, apperr.NewInvalidArgErr("invalid user id", err)
	}

	if err := s.repo.UpdateDeviceProfile(dp); err != nil {
		s.log.Error("device_profile.update failed: %v", err)
		return nil, mapRepoErr("update device profile", err)
	}
	return dp, nil
}

func (s *DeviceProfileServiceImpl) DeleteDeviceProfile(ctx context.Context, id string) error {
	s.log.Trace("device_profile.delete", "id", id)
	userID := ctx.Value(middleware.AuthUserIDKey).(string)

	if err := s.v.Var(id, "required,uuid4"); err != nil {
		return apperr.NewInvalidArgErr("invalid id", err)
	}

	if err := s.repo.DeleteDeviceProfile(userID, id); err != nil {
		s.log.Error("device_profile.delete failed: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.NewNotFoundErr("device profile not found", err)
		}
		return mapRepoErr("delete device profile", err)
	}
	return nil
}

func (s *DeviceProfileServiceImpl) createDeviceProfileByTemplate(t *entity.DeviceTemplate) *entity.DeviceProfile {
	var dp entity.DeviceProfile
	dp.Name = t.Name
	dp.DeviceType = t.DeviceType
	dp.Width = t.Width
	dp.Height = t.Height
	dp.TemplateID = nil
	dp.CustomHeaders = t.DefaultHeaders
	dp.UserAgent = &t.UserAgent
	dp.CountryCode = t.CountryCode
	dp.CreatedAt = time.Now()
	return &dp
}

func mapRepoErr(action string, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperr.NewNotFoundErr(action, err)
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return apperr.NewAlreadyExistsErr(action, err)
		case "23503":
			return apperr.NewInvalidArgErr(action, err)
		case "23514":
			return apperr.NewInvalidArgErr(action, err)
		default:
			return apperr.NewInternalErr(action, err)
		}
	}
	return apperr.NewInternalErr(action, err)
}
