package test

import (
	"testing"
	"time"

	"zenrows-challenge/internal/adapter/repo"
	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/test/util"

	"zenrows-challenge/internal/pkg/applog"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

func setupDPRepo(t *testing.T) (*repo.DeviceProfileRepoImpl, *entity.User) {
	t.Helper()
	require.NoError(t, util.LoadConfig())
	if _, err := util.InitTestContainers(t); err != nil {
		require.NoError(t, err)
	}
	dbConn, err := util.NewTestDB()
	require.NoError(t, err)
	logger := applog.NewAppDefaultLogger()
	r := repo.NewDeviceProfileRepoImpl(logger, dbConn)

	// Unique user for isolation
	uname := "testuser_" + uuid.NewString()
	passHash, _ := bcrypt.GenerateFromPassword([]byte("pass1234"), bcrypt.DefaultCost)
	u := entity.User{Username: uname, PasswordHash: string(passHash)}
	require.NoError(t, dbConn.Create(&u).Error)
	require.NotZero(t, u.ID)
	return r, &u
}

func TestDeviceProfileRepo_CreateDeviceProfile(t *testing.T) {
	r, u := setupDPRepo(t)

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "creates profile with optional fields",
			run: func(t *testing.T) {
				width := 1024
				height := 768
				ua := "Mozilla/5.0"
				cc := "US"
				dp := entity.DeviceProfile{
					UserID:        u.ID,
					Name:          "Profile A",
					DeviceType:    "desktop",
					Width:         &width,
					Height:        &height,
					UserAgent:     &ua,
					CountryCode:   &cc,
					CustomHeaders: datatypes.JSONMap{"X-Test": "true"},
				}
				require.NoError(t, r.CreateDeviceProfile(&dp))
				assert.NotZero(t, dp.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func TestDeviceProfileRepo_ListDeviceProfilesRetrievesCreatedProfile(t *testing.T) {
	r, u := setupDPRepo(t)
	dp := entity.DeviceProfile{UserID: u.ID, Name: "P1", DeviceType: "desktop"}
	require.NoError(t, r.CreateDeviceProfile(&dp))

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "fetches by id",
			run: func(t *testing.T) {
				items, err := r.ListDeviceProfiles(u.ID.String(), 1, 10)
				require.NoError(t, err)
				var got *entity.DeviceProfile
				for i := range items {
					if items[i].ID == dp.ID {
						got = &items[i]
						break
					}
				}
				require.NotNil(t, got)
				assert.Equal(t, dp.ID, got.ID)
				assert.Equal(t, u.ID, got.UserID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func TestDeviceProfileRepo_UpdateDeviceProfileSelective(t *testing.T) {
	r, u := setupDPRepo(t)
	width := 800
	ua := "UA/0"
	dp := entity.DeviceProfile{UserID: u.ID, Name: "Psel", DeviceType: "desktop", Width: &width, UserAgent: &ua}
	require.NoError(t, r.CreateDeviceProfile(&dp))

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "updates only provided non-zero fields",
			run: func(t *testing.T) {
				newWidth := 1280
				newUA := "UA/1"
				patch := entity.DeviceProfile{ID: dp.ID, UserID: dp.UserID, Name: "Psel2", Width: &newWidth, UserAgent: &newUA}
				require.NoError(t, r.UpdateDeviceProfile(&patch))

				items, err := r.ListDeviceProfiles(u.ID.String(), 1, 10)
				require.NoError(t, err)
				var got *entity.DeviceProfile
				for i := range items {
					if items[i].ID == dp.ID {
						got = &items[i]
						break
					}
				}
				require.NotNil(t, got)
				assert.Equal(t, "Psel2", got.Name)
				if assert.NotNil(t, got.Width) {
					assert.Equal(t, newWidth, *got.Width)
				}
				if assert.NotNil(t, got.UserAgent) {
					assert.Equal(t, newUA, *got.UserAgent)
				}
				assert.Equal(t, "desktop", got.DeviceType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func TestDeviceProfileRepo_ListDeviceProfiles(t *testing.T) {
	r, u := setupDPRepo(t)
	// Create multiple profiles
	for _, n := range []string{"L1", "L2", "L3"} {
		tmp := entity.DeviceProfile{UserID: u.ID, Name: n, DeviceType: "desktop"}
		require.NoError(t, r.CreateDeviceProfile(&tmp))
		time.Sleep(5 * time.Millisecond)
	}

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "paginates results",
			run: func(t *testing.T) {
				p1, err := r.ListDeviceProfiles(u.ID.String(), 1, 2)
				require.NoError(t, err)
				assert.Len(t, p1, 2)

				p2, err := r.ListDeviceProfiles(u.ID.String(), 2, 2)
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(p2), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func TestDeviceProfileRepo_UpdateDeviceProfile(t *testing.T) {
	r, u := setupDPRepo(t)
	dp := entity.DeviceProfile{UserID: u.ID, Name: "PU1", DeviceType: "desktop"}
	require.NoError(t, r.CreateDeviceProfile(&dp))

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "updates full entity",
			run: func(t *testing.T) {
				dp.DeviceType = "mobile"
				require.NoError(t, r.UpdateDeviceProfile(&dp))
				items, err := r.ListDeviceProfiles(u.ID.String(), 1, 10)
				require.NoError(t, err)
				var got *entity.DeviceProfile
				for i := range items {
					if items[i].ID == dp.ID {
						got = &items[i]
						break
					}
				}
				require.NotNil(t, got)
				assert.Equal(t, "mobile", got.DeviceType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func TestDeviceProfileRepo_DeleteDeviceProfile(t *testing.T) {
	r, u := setupDPRepo(t)
	dp := entity.DeviceProfile{UserID: u.ID, Name: "PD1", DeviceType: "desktop"}
	require.NoError(t, r.CreateDeviceProfile(&dp))

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "deletes by id",
			run: func(t *testing.T) {
				require.NoError(t, r.DeleteDeviceProfile(u.ID.String(), dp.ID.String()))
				items, err := r.ListDeviceProfiles(u.ID.String(), 1, 10)
				require.NoError(t, err)
				assert.Len(t, items, 0)
				for _, item := range items {
					require.NotEqual(t, dp.ID, item.ID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
