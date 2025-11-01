package repo

import (
	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/pkg/applog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceTemplateRepoImpl struct {
	log applog.AppLogger
	db  *gorm.DB
}

func NewDeviceTemplateRepoImpl(log applog.AppLogger, db *gorm.DB) *DeviceTemplateRepoImpl {
	return &DeviceTemplateRepoImpl{log: log, db: db}
}

func (r *DeviceTemplateRepoImpl) GetDeviceTemplates() ([]entity.DeviceTemplate, error) {
	r.log.Trace("device_template.list")
	var out []entity.DeviceTemplate
	if err := r.db.Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *DeviceTemplateRepoImpl) GetDeviceTemplateByID(id *uuid.UUID) (*entity.DeviceTemplate, error) {
	r.log.Trace("device_template.get", "id", id.String())
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	r.log.Trace("device_template.get", "id", id.String())
	var out entity.DeviceTemplate
	if err := r.db.First(&out, "id = ?", *id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
