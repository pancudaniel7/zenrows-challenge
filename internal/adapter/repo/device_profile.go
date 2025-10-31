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

// ListDeviceProfiles returns device profiles for a given user ID using pagination.
// Page numbering starts at 1. If page < 1 or pageSize <= 0, sensible defaults are applied.
func (r *DeviceProfileRepoImpl) ListDeviceProfiles(userID string, page, pageSize int) ([]entity.DeviceProfile, error) {
	r.log.Trace("device_profile.list", "user_id", userID)

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
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

func (r *DeviceProfileRepoImpl) UpdateDeviceProfileSelective(dp *entity.DeviceProfile) error {
	r.log.Trace("device_profile.update_selective", "id", dp.ID.String(), "user_id", dp.UserID.String())
	return r.db.Model(&entity.DeviceProfile{}).
		Where("id = ? AND user_id = ?", dp.ID, dp.UserID).
		Updates(dp).Error
}

func (r *DeviceProfileRepoImpl) GetDeviceProfileByID(id string) (*entity.DeviceProfile, error) {
	r.log.Trace("device_profile.get", "id", id)
	var out entity.DeviceProfile
	if err := r.db.Where("id = ?", id).First(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *DeviceProfileRepoImpl) UpdateDeviceProfile(dp *entity.DeviceProfile) error {
	r.log.Trace("device_profile.update", "id", dp.ID.String(), "name", dp.Name)
	return r.db.Save(dp).Error
}

func (r *DeviceProfileRepoImpl) DeleteDeviceProfile(id string) error {
    r.log.Trace("device_profile.delete", "id", id)
    pid, err := uuid.Parse(id)
    if err != nil {
        return err
    }
    return r.db.Where("id = ?", pid).Delete(&entity.DeviceProfile{}).Error
}
