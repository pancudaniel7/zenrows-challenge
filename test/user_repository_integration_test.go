package test

import (
	"testing"
	"zenrows-challenge/internal/adapter/repo"
	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/test/util"

	"zenrows-challenge/internal/pkg/applog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepositoryIntegration(t *testing.T) {
	if err := util.LoadConfig(); err != nil {
		assert.FailNow(t, err.Error())
	}

	dbConn, err := util.NewTestDB()
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	logger := applog.NewAppDefaultLogger()
	r := repo.NewUserRepoImpl(logger, dbConn)

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "check user credentials",
			run: func(t *testing.T) {
				passwordHash, err := bcrypt.GenerateFromPassword([]byte("alicepass"), bcrypt.DefaultCost)
				if err != nil {
					assert.FailNow(t, err.Error())
				}
				user := entity.User{Username: "alice", PasswordHash: string(passwordHash)}
				userID, passHash, err := r.RetrieveCredentials(user)
				if err != nil {
					return
				}

				require.NoError(t, err)
				require.NotZero(t, userID)
				require.NotZero(t, passHash)
			},
		},
		{
			name: "check user wrong credentials",
			run: func(t *testing.T) {
				passwordHash, err := bcrypt.GenerateFromPassword([]byte("alicewrongpass"), bcrypt.DefaultCost)
				if err != nil {
					assert.FailNow(t, err.Error())
				}
				user := entity.User{Username: "wrong", PasswordHash: string(passwordHash)}
				userID, passHash, err := r.RetrieveCredentials(user)

				assert.NoError(t, err)
				assert.Zero(t, userID)
				assert.Zero(t, passHash)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
