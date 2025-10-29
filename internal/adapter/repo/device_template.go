package repo

import (
	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/pkg/applog"

	"gorm.io/gorm"
)

type DeviceTemplateImplementation struct {
	log applog.AppLogger
	db  *gorm.DB
}

func NewDeviceTemplateImplementation(log applog.AppLogger, db *gorm.DB) *DeviceTemplateImplementation {
	return &DeviceTemplateImplementation{log: log, db: db}
}

func (r *DeviceTemplateImplementation) GetDeviceTemplates() ([]entity.DeviceTemplate, error) {
	var out []entity.DeviceTemplate
	if err := r.db.Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
