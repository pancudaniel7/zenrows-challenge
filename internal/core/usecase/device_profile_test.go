package usecase

import (
	"context"
	"errors"
	"testing"

	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/pkg/apperr"
	"zenrows-challenge/internal/pkg/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type mockDeviceProfileRepo struct {
	createFn func(*entity.DeviceProfile) error
	listFn   func(string, int, int) ([]entity.DeviceProfile, error)
	updateFn func(*entity.DeviceProfile) error
	deleteFn func(string, string) error
}

func (m *mockDeviceProfileRepo) CreateDeviceProfile(dp *entity.DeviceProfile) error {
	if m.createFn != nil {
		return m.createFn(dp)
	}
	return nil
}

func (m *mockDeviceProfileRepo) ListDeviceProfiles(userID string, page, pageSize int) ([]entity.DeviceProfile, error) {
	if m.listFn != nil {
		return m.listFn(userID, page, pageSize)
	}
	return nil, nil
}

func (m *mockDeviceProfileRepo) UpdateDeviceProfile(dp *entity.DeviceProfile) error {
	if m.updateFn != nil {
		return m.updateFn(dp)
	}
	return nil
}

func (m *mockDeviceProfileRepo) UpdateDeviceProfileSelective(dp *entity.DeviceProfile) error {
	return m.UpdateDeviceProfile(dp)
}

func (m *mockDeviceProfileRepo) DeleteDeviceProfile(userID, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(userID, id)
	}
	return nil
}

type mockDeviceTemplateRepo struct {
	getFn func(*uuid.UUID) (*entity.DeviceTemplate, error)
}

func (m *mockDeviceTemplateRepo) GetDeviceTemplates() ([]entity.DeviceTemplate, error) {
	return nil, errors.New("not implemented")
}

func (m *mockDeviceTemplateRepo) GetDeviceTemplateByID(id *uuid.UUID) (*entity.DeviceTemplate, error) {
	if m.getFn != nil {
		return m.getFn(id)
	}
	return nil, nil
}

func TestDeviceProfileService_CreateDeviceProfile_Success(t *testing.T) {
	repoCalled := false
	repo := &mockDeviceProfileRepo{
		createFn: func(dp *entity.DeviceProfile) error {
			repoCalled = true
			assert.Equal(t, "Profile", dp.Name)
			return nil
		},
	}
	svc := NewDeviceProfileServiceImpl(noopLogger{}, repo, &mockDeviceTemplateRepo{}, validator.New())
	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, uuid.NewString())

	dp := &entity.DeviceProfile{
		UserID:     uuid.New(),
		Name:       "Profile",
		DeviceType: "desktop",
	}

	err := svc.CreateDeviceProfile(ctx, dp)
	require.NoError(t, err)
	assert.True(t, repoCalled)
}

func TestDeviceProfileService_CreateDeviceProfile_InvalidPayload(t *testing.T) {
	svc := NewDeviceProfileServiceImpl(noopLogger{}, &mockDeviceProfileRepo{}, &mockDeviceTemplateRepo{}, validator.New())
	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, uuid.NewString())
	dp := &entity.DeviceProfile{DeviceType: "desktop"}

	err := svc.CreateDeviceProfile(ctx, dp)
	require.Error(t, err)
	var inv *apperr.InvalidArgErr
	assert.ErrorAs(t, err, &inv)
}

func TestDeviceProfileService_CreateDeviceProfile_UsesTemplate(t *testing.T) {
	templateID := uuid.New()
	userID := uuid.New()
	templateRepoCalled := false
	repoCalled := false

	templateRepo := &mockDeviceTemplateRepo{
		getFn: func(id *uuid.UUID) (*entity.DeviceTemplate, error) {
			templateRepoCalled = true
			assert.Equal(t, templateID, *id)
			return &entity.DeviceTemplate{
				Name:           "Template",
				DeviceType:     "mobile",
				UserAgent:      "UA",
				DefaultHeaders: datatypes.JSONMap{"X-Test": "true"},
			}, nil
		},
	}
	repo := &mockDeviceProfileRepo{
		createFn: func(dp *entity.DeviceProfile) error {
			repoCalled = true
			assert.Equal(t, "Template", dp.Name)
			assert.Equal(t, userID, dp.UserID)
			assert.Equal(t, "mobile", dp.DeviceType)
			return nil
		},
	}
	svc := NewDeviceProfileServiceImpl(noopLogger{}, repo, templateRepo, validator.New())

	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, userID.String())
	dp := &entity.DeviceProfile{
		UserID:     uuid.New(),
		Name:       "ignored",
		DeviceType: "desktop",
		TemplateID: &templateID,
	}

	err := svc.CreateDeviceProfile(ctx, dp)
	require.NoError(t, err)
	assert.True(t, templateRepoCalled)
	assert.True(t, repoCalled)
}

func TestDeviceProfileService_CreateDeviceProfile_RepoError(t *testing.T) {
	repo := &mockDeviceProfileRepo{
		createFn: func(*entity.DeviceProfile) error {
			return &pgconn.PgError{Code: "23505"}
		},
	}
	svc := NewDeviceProfileServiceImpl(noopLogger{}, repo, &mockDeviceTemplateRepo{}, validator.New())
	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, uuid.NewString())

	dp := &entity.DeviceProfile{
		UserID:     uuid.New(),
		Name:       "dup",
		DeviceType: "desktop",
	}

	err := svc.CreateDeviceProfile(ctx, dp)
	require.Error(t, err)
	var ae *apperr.AlreadyExistsErr
	assert.ErrorAs(t, err, &ae)
}

func TestDeviceProfileService_ListDeviceProfilesByUserID(t *testing.T) {
	repo := &mockDeviceProfileRepo{
		listFn: func(userID string, page, pageSize int) ([]entity.DeviceProfile, error) {
			assert.Equal(t, "user", userID)
			assert.Equal(t, 1, page)
			assert.Equal(t, 10, pageSize)
			return []entity.DeviceProfile{{Name: "A"}}, nil
		},
	}
	svc := NewDeviceProfileServiceImpl(noopLogger{}, repo, &mockDeviceTemplateRepo{}, validator.New())

	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, "user")
	out, err := svc.ListDeviceProfilesByUserID(ctx, 1, 10)
	require.NoError(t, err)
	assert.Len(t, out, 1)
}

func TestDeviceProfileService_UpdateDeviceProfile_NotAuthorized(t *testing.T) {
	svc := NewDeviceProfileServiceImpl(noopLogger{}, &mockDeviceProfileRepo{}, &mockDeviceTemplateRepo{}, validator.New())
	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, "some-user")
	dp := &entity.DeviceProfile{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		Name:       "Profile",
		DeviceType: "desktop",
	}

	_, err := svc.UpdateDeviceProfile(ctx, dp)
	require.Error(t, err)
	var na *apperr.NotAuthorizedErr
	assert.ErrorAs(t, err, &na)
}

func TestDeviceProfileService_UpdateDeviceProfile_Success(t *testing.T) {
	repoCalled := false
	repo := &mockDeviceProfileRepo{
		updateFn: func(dp *entity.DeviceProfile) error {
			repoCalled = true
			assert.Equal(t, "Updated", dp.Name)
			return nil
		},
	}
	svc := NewDeviceProfileServiceImpl(noopLogger{}, repo, &mockDeviceTemplateRepo{}, validator.New())

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, userID.String())
	dp := &entity.DeviceProfile{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       "Updated",
		DeviceType: "desktop",
	}

	out, err := svc.UpdateDeviceProfile(ctx, dp)
	require.NoError(t, err)
	assert.True(t, repoCalled)
	assert.Equal(t, "Updated", out.Name)
}

func TestDeviceProfileService_UpdateDeviceProfile_MapError(t *testing.T) {
	repo := &mockDeviceProfileRepo{
		updateFn: func(*entity.DeviceProfile) error {
			return gorm.ErrRecordNotFound
		},
	}
	svc := NewDeviceProfileServiceImpl(noopLogger{}, repo, &mockDeviceTemplateRepo{}, validator.New())

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, userID.String())
	dp := &entity.DeviceProfile{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       "Name",
		DeviceType: "desktop",
	}

	_, err := svc.UpdateDeviceProfile(ctx, dp)
	require.Error(t, err)
	var nf *apperr.NotFoundErr
	assert.ErrorAs(t, err, &nf)
}

func TestDeviceProfileService_DeleteDeviceProfile_InvalidID(t *testing.T) {
	svc := NewDeviceProfileServiceImpl(noopLogger{}, &mockDeviceProfileRepo{}, &mockDeviceTemplateRepo{}, validator.New())
	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, uuid.NewString())

	err := svc.DeleteDeviceProfile(ctx, "not-a-uuid")
	require.Error(t, err)
	var inv *apperr.InvalidArgErr
	assert.ErrorAs(t, err, &inv)
}

func TestDeviceProfileService_DeleteDeviceProfile_Success(t *testing.T) {
	repoCalled := false
	repo := &mockDeviceProfileRepo{
		deleteFn: func(userID, id string) error {
			repoCalled = true
			assert.Equal(t, "user", userID)
			return nil
		},
	}
	svc := NewDeviceProfileServiceImpl(noopLogger{}, repo, &mockDeviceTemplateRepo{}, validator.New())
	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, "user")

	err := svc.DeleteDeviceProfile(ctx, uuid.NewString())
	require.NoError(t, err)
	assert.True(t, repoCalled)
}

func TestDeviceProfileService_DeleteDeviceProfile_NotFound(t *testing.T) {
	repo := &mockDeviceProfileRepo{
		deleteFn: func(string, string) error {
			return gorm.ErrRecordNotFound
		},
	}
	svc := NewDeviceProfileServiceImpl(noopLogger{}, repo, &mockDeviceTemplateRepo{}, validator.New())
	ctx := context.WithValue(context.Background(), middleware.AuthUserIDKey, uuid.NewString())

	err := svc.DeleteDeviceProfile(ctx, uuid.NewString())
	require.Error(t, err)
	var nf *apperr.NotFoundErr
	assert.ErrorAs(t, err, &nf)
}
