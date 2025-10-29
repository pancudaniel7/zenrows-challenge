package usecase

import (
	"sync"
	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/core/port"
	"zenrows-challenge/internal/pkg/apperr"
	"zenrows-challenge/internal/pkg/applog"

	"golang.org/x/crypto/bcrypt"
)

type AuthenticationServiceImpl struct {
	log      applog.AppLogger
	wg       *sync.WaitGroup
	userRepo port.UserRepo
}

func NewAuthenticationService(log applog.AppLogger, wg *sync.WaitGroup, ur port.UserRepo) *AuthenticationServiceImpl {
	return &AuthenticationServiceImpl{log: log, wg: wg, userRepo: ur}
}

func (s *AuthenticationServiceImpl) CheckCredentials(username string, password string) (string, error) {
	user := entity.User{
		Username: username,
	}

	userID, passwordHash, err := s.userRepo.RetrieveCredentials(user)
	if err != nil {
		return "", &apperr.NotAuthorizedErr{
			Msg:   "Unauthorized",
			Cause: err,
		}
	}

	if userID == "" || passwordHash == "" {
		return "", &apperr.NotAuthorizedErr{
			Msg:   "Unauthorized",
			Cause: err,
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil {
		return "", nil
	}
	return userID, nil
}
