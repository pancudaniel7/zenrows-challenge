package port

import "zenrows-challenge/internal/core/entity"

type UserRepo interface {
	RetrieveCredentials(u entity.User) (string, string, error)
}
