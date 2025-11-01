package usecase

import (
	"errors"
	"testing"

	"zenrows-challenge/internal/core/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type deviceTemplateRepoMock struct {
	templates []entity.DeviceTemplate
	err       error
	called    bool
}

func (m *deviceTemplateRepoMock) GetDeviceTemplates() ([]entity.DeviceTemplate, error) {
	m.called = true
	if m.err != nil {
		return nil, m.err
	}
	return m.templates, nil
}

func (m *deviceTemplateRepoMock) GetDeviceTemplateByID(_ *uuid.UUID) (*entity.DeviceTemplate, error) {
	return nil, errors.New("not implemented")
}

// Ensure interface compliance.
var _ interface {
	GetDeviceTemplates() ([]entity.DeviceTemplate, error)
	GetDeviceTemplateByID(*uuid.UUID) (*entity.DeviceTemplate, error)
} = (*deviceTemplateRepoMock)(nil)

type noopLogger struct{}

func (noopLogger) Info(string, ...any)  {}
func (noopLogger) Warn(string, ...any)  {}
func (noopLogger) Error(string, ...any) {}
func (noopLogger) Debug(string, ...any) {}
func (noopLogger) Trace(string, ...any) {}
func (noopLogger) Fatal(string, ...any) {}

func TestDeviceTemplateService_RetrieveDeviceTemplates(t *testing.T) {
	want := []entity.DeviceTemplate{{Name: "Desktop"}}
	repo := &deviceTemplateRepoMock{templates: want}
	svc := NewDeviceTemplateServiceImpl(noopLogger{}, repo)

	got, err := svc.RetrieveDeviceTemplates()
	require.NoError(t, err)
	assert.True(t, repo.called, "expected repo to be called")
	assert.Equal(t, want, got)
}

func TestDeviceTemplateService_RetrieveDeviceTemplatesError(t *testing.T) {
	repo := &deviceTemplateRepoMock{err: errors.New("boom")}
	svc := NewDeviceTemplateServiceImpl(noopLogger{}, repo)

	_, err := svc.RetrieveDeviceTemplates()
	require.Error(t, err)
	assert.True(t, repo.called, "expected repo to be called")
}
