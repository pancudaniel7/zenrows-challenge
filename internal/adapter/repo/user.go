package repo

import (
	"errors"

	"zenrows-challenge/internal/core/entity"

	"gorm.io/gorm"
)

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepoImpl(db *gorm.DB) *UserRepoImpl { return &UserRepoImpl{db: db} }

func (r *UserRepoImpl) RetrieveCredentials(u entity.User) (string, string, error) {
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
