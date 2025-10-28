package repo

import (
	"errors"

	"zenrows-challenge/internal/core/entity"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) RetrieveCredentials(u entity.User) (string, string, error) {
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
