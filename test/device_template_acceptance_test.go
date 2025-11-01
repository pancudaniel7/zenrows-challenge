package test

import (
	"encoding/json"
	"io"
	"net"
	nethttp "net/http"
	"testing"
	"time"

	httpadapter "zenrows-challenge/internal/adapter/http"
	"zenrows-challenge/internal/adapter/repo"
	"zenrows-challenge/internal/core/entity"
	"zenrows-challenge/internal/core/port"
	"zenrows-challenge/internal/core/usecase"
	"zenrows-challenge/internal/pkg/apperr"
	"zenrows-challenge/internal/pkg/middleware"
	testutil "zenrows-challenge/test/util"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type templateSuite struct {
	app     *fiber.App
	client  *nethttp.Client
	baseURL string
	db      *gorm.DB
	repo    *repo.DeviceTemplateRepoImpl
	cleanup func()
}

func (s *templateSuite) truncateTemplates(t *testing.T) {
	t.Helper()
	if s.db == nil {
		return
	}
	require.NoError(t, s.db.Exec("TRUNCATE TABLE zenrows.device_template RESTART IDENTITY CASCADE").Error)
}

func (s *templateSuite) doGet(t *testing.T, path string, headers map[string]string) *nethttp.Response {
	t.Helper()
	req, err := nethttp.NewRequest(nethttp.MethodGet, s.baseURL+path, nil)
	require.NoError(t, err)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := s.client.Do(req)
	require.NoError(t, err)
	return resp
}

func newDeviceTemplateSuite(t *testing.T, withAuth bool, svc port.DeviceTemplateService) *templateSuite {
	t.Helper()

	require.NoError(t, testutil.LoadConfig())

	logger := noopLogger{}
	var (
		dbConn   *gorm.DB
		err      error
		repoImpl *repo.DeviceTemplateRepoImpl
	)

	if svc == nil {
		t.Setenv("SKIP_DEVICE_TEMPLATE_SEED", "true")
		if _, err = testutil.InitTestContainers(t); err != nil {
			require.NoError(t, err)
		}
		dbConn, err = testutil.NewTestDB()
		require.NoError(t, err)

		repoImpl = repo.NewDeviceTemplateRepoImpl(logger, dbConn)
		svc = usecase.NewDeviceTemplateServiceImpl(logger, repoImpl)
	}

	handler := httpadapter.NewDeviceTemplateHandlerImpl(logger, svc)

	app := fiber.New()
	if withAuth {
		app.Use(func(c fiber.Ctx) error {
			if c.Get("Authorization") != basicAuthHeader {
				return c.Status(nethttp.StatusUnauthorized).JSON(map[string]string{
					"code":    "NOT_AUTHORIZED",
					"message": "unauthorized",
				})
			}
			c.Locals(middleware.AuthUserIDKey, uuid.New().String())
			return c.Next()
		})
	}
	app.Get("/device-templates", handler.List)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() {
		if err := app.Listener(listener); err != nil && err != net.ErrClosed {
			panic(err)
		}
	}()

	baseURL := "http://" + listener.Addr().String()
	client := &nethttp.Client{Timeout: 5 * time.Second}
	waitForServer(t, client, baseURL+"/device-templates")

	suite := &templateSuite{
		app:     app,
		client:  client,
		baseURL: baseURL,
		db:      dbConn,
		repo:    repoImpl,
		cleanup: func() {
			_ = app.Shutdown()
			_ = listener.Close()
		},
	}
	suite.truncateTemplates(t)
	t.Cleanup(suite.cleanup)
	return suite
}

type erroringTemplateService struct {
	err error
}

func (e erroringTemplateService) RetrieveDeviceTemplates() ([]entity.DeviceTemplate, error) {
	return nil, e.err
}

func TestDeviceTemplateHandler_List(t *testing.T) {
	cases := []struct {
		name           string
		path           string
		withAuth       bool
		headers        map[string]string
		setup          func(t *testing.T, suite *templateSuite)
		svc            port.DeviceTemplateService
		expectedStatus int
		assertFn       func(t *testing.T, status int, payload []byte)
	}{
		{
			name:     "returns templates",
				path:     "/device-templates",
				withAuth: true,
				setup: func(t *testing.T, suite *templateSuite) {
					suite.truncateTemplates(t)
					require.NotNil(t, suite.db)
					width := 1024
					cc := "US"
					tpl := entity.DeviceTemplate{
						Name:           "Desktop",
					DeviceType:     "desktop",
					Width:          &width,
					UserAgent:      "Mozilla/5.0",
					CountryCode:    &cc,
					DefaultHeaders: datatypes.JSONMap{"X-Test": "true"},
				}
				require.NoError(t, suite.db.Create(&tpl).Error)
			},
			expectedStatus: nethttp.StatusOK,
			assertFn: func(t *testing.T, status int, payload []byte) {
				require.Equal(t, nethttp.StatusOK, status)
				var out []entity.DeviceTemplate
				require.NoError(t, json.Unmarshal(payload, &out))
				require.Len(t, out, 1)
				assert.Equal(t, "Desktop", out[0].Name)
			},
		},
		{
			name:           "returns empty list when no templates",
			path:           "/device-templates",
			withAuth:       true,
			setup: func(t *testing.T, suite *templateSuite) {
				suite.truncateTemplates(t)
			},
			expectedStatus: nethttp.StatusOK,
			assertFn: func(t *testing.T, status int, payload []byte) {
				require.Equal(t, nethttp.StatusOK, status)
				var out []entity.DeviceTemplate
				require.NoError(t, json.Unmarshal(payload, &out))
				assert.Len(t, out, 0)
			},
		},
		{
			name:           "missing auth returns unauthorized",
			path:           "/device-templates",
			withAuth:       true,
			headers:        map[string]string{},
			expectedStatus: nethttp.StatusUnauthorized,
			assertFn: func(t *testing.T, status int, payload []byte) {
				require.Equal(t, nethttp.StatusUnauthorized, status)
				var out map[string]string
				require.NoError(t, json.Unmarshal(payload, &out))
				assert.Equal(t, "NOT_AUTHORIZED", out["code"])
			},
		},
		{
			name:           "maps domain error to http",
			path:           "/device-templates",
			withAuth:       true,
			svc:            erroringTemplateService{err: apperr.NewInvalidArgErr("bad request", nil)},
			expectedStatus: nethttp.StatusBadRequest,
			assertFn: func(t *testing.T, status int, payload []byte) {
				require.Equal(t, nethttp.StatusBadRequest, status)
				var out map[string]string
				require.NoError(t, json.Unmarshal(payload, &out))
				assert.Equal(t, "INVALID_ARGUMENT", out["code"])
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			suite := newDeviceTemplateSuite(t, tc.withAuth, tc.svc)
			if tc.setup != nil {
				tc.setup(t, suite)
			}

			headers := map[string]string{}
			if tc.headers != nil {
				for k, v := range tc.headers {
					headers[k] = v
				}
			} else if tc.withAuth {
				headers["Authorization"] = basicAuthHeader
			}

			resp := suite.doGet(t, tc.path, headers)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, resp.StatusCode)
			tc.assertFn(t, resp.StatusCode, body)
		})
	}
}
