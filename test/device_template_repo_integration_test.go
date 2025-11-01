package test

import (
	"testing"
	"zenrows-challenge/internal/adapter/repo"
	"zenrows-challenge/test/util"

	"zenrows-challenge/internal/pkg/applog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceTemplateRepositoryIntegration(t *testing.T) {
	if err := util.LoadConfig(); err != nil {
		assert.FailNow(t, err.Error())
	}

	if _, err := util.InitTestContainers(t); err != nil {
		assert.FailNow(t, err.Error())
	}

	dbConn, err := util.NewTestDB()
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	logger := applog.NewAppDefaultLogger()
	r := repo.NewDeviceTemplateRepoImpl(logger, dbConn)

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "retrieve device templates",
			run: func(t *testing.T) {

				dts, err := r.GetDeviceTemplates()
				require.NoError(t, err)

				for i, dt := range dts {
					t.Logf("device template[%d]: id=%s name=%s device_type=%s", i, dt.ID.String(), dt.Name, dt.DeviceType)
					assert.NotZero(t, dt.ID, "template ID should not be zero")
					assert.NotEmpty(t, dt.Name, "template Name should not be empty")
					assert.Contains(t, []string{"desktop", "mobile"}, dt.DeviceType, "device type should be desktop or mobile")
				}

				if len(dts) == 0 {
					t.Log("no device templates found in test DB")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
