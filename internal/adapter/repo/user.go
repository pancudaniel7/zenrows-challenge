package repo

import (
	"errors"

	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/pkg/applog"

	"gorm.io/gorm"
)

type UserRepoImpl struct {
	log applog.AppLogger
	db  *gorm.DB
}

func NewUserRepoImpl(log applog.AppLogger, db *gorm.DB) *UserRepoImpl {
	return &UserRepoImpl{log: log, db: db}
}

func (r *UserRepoImpl) RetrieveCredentials(u entity.User) (string, string, error) {
	r.log.Trace("user.retrieve_credentials", "username", u.Username)
	var found entity.User

	err := r.db.Where("username = ?", u.Username).First(&found).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil
		}
		return "", "", err
	}
	return found.ID.String(), found.PasswordHash, nil
}
