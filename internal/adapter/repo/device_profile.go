package repo

import (
	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/pkg/applog"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DeviceProfileRepoImpl struct {
	log applog.AppLogger
	db  *gorm.DB
}

func NewDeviceProfileRepoImpl(log applog.AppLogger, db *gorm.DB) *DeviceProfileRepoImpl {
	return &DeviceProfileRepoImpl{log: log, db: db}
}

func (r *DeviceProfileRepoImpl) ListDeviceProfiles(userID string, page, pageSize int) ([]entity.DeviceProfile, error) {
	r.log.Trace("device_profile.list", "user_id", userID)

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	var out []entity.DeviceProfile
	if err := r.db.Where("user_id = ?", uid).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *DeviceProfileRepoImpl) CreateDeviceProfile(dp *entity.DeviceProfile) error {
	if dp.CustomHeaders == nil {
		dp.CustomHeaders = datatypes.JSONMap{}
	}
	r.log.Trace("device_profile.create", "user_id", dp.UserID.String(), "name", dp.Name)
	return r.db.Create(dp).Error
}

func (r *DeviceProfileRepoImpl) UpdateDeviceProfile(dp *entity.DeviceProfile) error {
	r.log.Trace("device_profile.update_selective", "id", dp.ID.String(), "user_id", dp.UserID.String())
	return r.db.Model(&entity.DeviceProfile{}).
		Where("id = ? AND user_id = ?", dp.ID, dp.UserID).
		Updates(dp).Error
}

func (r *DeviceProfileRepoImpl) DeleteDeviceProfile(userID, id string) error {
	r.log.Trace("device_profile.delete", "id", id)
	pid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Where("id = ? AND user_id = ?", pid, userID).Delete(&entity.DeviceProfile{}).Error
}
