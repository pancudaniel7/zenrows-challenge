package usecase

import (
	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/core/port"
	"zenrows-challenge/internal/pkg/applog"
)

type DeviceTemplateServiceImpl struct {
	repo port.DeviceTemplateRepo
	log  applog.AppLogger
}

func NewDeviceTemplateServiceImpl(log applog.AppLogger, r port.DeviceTemplateRepo) *DeviceTemplateServiceImpl {
	return &DeviceTemplateServiceImpl{repo: r, log: log}
}

func (s *DeviceTemplateServiceImpl) RetrieveDeviceTemplates() ([]entity.DeviceTemplate, error) {
	s.log.Trace("DeviceTemplateServiceImpl: RetrieveDeviceTemplates called")
	return s.repo.GetDeviceTemplates()
}
