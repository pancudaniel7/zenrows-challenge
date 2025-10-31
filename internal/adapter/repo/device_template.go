package repo

import (
	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/pkg/applog"

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
