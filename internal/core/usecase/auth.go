package usecase

import (
	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/core/port"
	"zenrows-challenge/internal/pkg/apperr"
	"zenrows-challenge/internal/pkg/applog"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationServiceImpl struct {
	log      applog.AppLogger
	userRepo port.UserRepo
}

func NewAuthenticationService(log applog.AppLogger, ur port.UserRepo) *AuthenticationServiceImpl {
	return &AuthenticationServiceImpl{log: log, userRepo: ur}
}

func (s *AuthenticationServiceImpl) CheckCredentials(username string, password string) (string, error) {
	user := entity.User{
		Username: username,
	}

	userID, passwordHash, err := s.userRepo.RetrieveCredentials(user)
	if err != nil {
		return "", apperr.NewNotAuthorizedErr("Unauthorized", err)
	}

	if userID == "" || passwordHash == "" {
		return "", apperr.NewNotAuthorizedErr("Unauthorized", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil {
		return "", nil
	}
	return userID, nil
}
