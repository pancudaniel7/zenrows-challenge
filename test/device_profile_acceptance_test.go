package test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type noopLogger struct{}

func (noopLogger) Info(string, ...any)  {}
func (noopLogger) Warn(string, ...any)  {}
func (noopLogger) Error(string, ...any) {}
func (noopLogger) Debug(string, ...any) {}
func (noopLogger) Trace(string, ...any) {}
func (noopLogger) Fatal(string, ...any) {}

var basicAuthHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))

type acceptanceSuite struct {
	app     *fiber.App
	client  *nethttp.Client
	baseURL string
	repo    *repo.DeviceProfileRepoImpl
	userID  uuid.UUID
	handler *httpadapter.DeviceProfileHandlerImpl
	cleanup func()
}

func (s *acceptanceSuite) doGet(t *testing.T, path string, headers map[string]string) *nethttp.Response {
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

func (s *acceptanceSuite) doPost(t *testing.T, path string, headers map[string]string, body []byte) *nethttp.Response {
	return s.doRequestWithBody(t, nethttp.MethodPost, path, headers, body)
}

func (s *acceptanceSuite) doPut(t *testing.T, path string, headers map[string]string, body []byte) *nethttp.Response {
	return s.doRequestWithBody(t, nethttp.MethodPut, path, headers, body)
}

func (s *acceptanceSuite) doDelete(t *testing.T, path string, headers map[string]string) *nethttp.Response {
	return s.doRequestWithBody(t, nethttp.MethodDelete, path, headers, nil)
}

func (s *acceptanceSuite) doRequestWithBody(t *testing.T, method, path string, headers map[string]string, body []byte) *nethttp.Response {
	t.Helper()
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := nethttp.NewRequest(method, s.baseURL+path, reader)
	require.NoError(t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := s.client.Do(req)
	require.NoError(t, err)
	return resp
}

func newDeviceProfileSuite(t *testing.T, withAuth bool, svc port.DeviceProfileService) *acceptanceSuite {
	t.Helper()

	require.NoError(t, testutil.LoadConfig())
	var (
		repository *repo.DeviceProfileRepoImpl
		dbConn     *gorm.DB
		err        error
	)

	logger := noopLogger{}
	v := validator.New()
	userID := uuid.New()

	if svc == nil {
		if _, err = testutil.InitTestContainers(t); err != nil {
			require.NoError(t, err)
		}
		dbConn, err = testutil.NewTestDB()
		require.NoError(t, err)

		repository = repo.NewDeviceProfileRepoImpl(logger, dbConn)
		svc = usecase.NewDeviceProfileServiceImpl(logger, repository, templateRepoStub{}, v)

		username := "accept_user_" + uuid.NewString()
		pw, err := bcrypt.GenerateFromPassword([]byte("pass1234"), bcrypt.DefaultCost)
		require.NoError(t, err)
		user := entity.User{Username: username, PasswordHash: string(pw)}
		require.NoError(t, dbConn.Create(&user).Error)
		userID = user.ID
	}

	handler := httpadapter.NewDeviceProfileHandlerImpl(logger, svc, v)

	app := fiber.New()
	if withAuth {
		uid := userID.String()
		app.Use(func(c fiber.Ctx) error {
			if c.Get("Authorization") != basicAuthHeader {
				return c.Status(nethttp.StatusUnauthorized).JSON(map[string]string{
					"code":    "NOT_AUTHORIZED",
					"message": "unauthorized",
				})
			}
			c.Locals(middleware.AuthUserIDKey, uid)
			return c.Next()
		})
	}
	app.Get("/device-profiles", handler.ListDeviceProfilesByUserID)
	app.Post("/device-profiles", handler.CreateDeviceProfile)
	app.Get("/device-profiles/:id", handler.GetDeviceProfileByID)
	app.Put("/device-profiles/:id", handler.UpdateDeviceProfile)
	app.Delete("/device-profiles/:id", handler.DeleteDeviceProfile)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() {
		if err := app.Listener(listener); err != nil && !errors.Is(err, net.ErrClosed) {
			panic(err)
		}
	}()

	baseURL := "http://" + listener.Addr().String()
	client := &nethttp.Client{Timeout: 5 * time.Second}
	waitForServer(t, client, baseURL+"/device-profiles")

	suite := &acceptanceSuite{
		app:     app,
		client:  client,
		baseURL: baseURL,
		repo:    repository,
		handler: handler,
		userID:  userID,
		cleanup: func() {
			_ = app.Shutdown()
			_ = listener.Close()
		},
	}
	t.Cleanup(suite.cleanup)
	return suite
}

type templateRepoStub struct{}

func (templateRepoStub) GetDeviceTemplates() ([]entity.DeviceTemplate, error) {
	return nil, gorm.ErrRecordNotFound
}

func (templateRepoStub) GetDeviceTemplateByID(id *uuid.UUID) (*entity.DeviceTemplate, error) {
	return nil, gorm.ErrRecordNotFound
}

func waitForServer(t *testing.T, client *nethttp.Client, url string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("server at %s did not start on time", url)
}

func TestListDeviceProfilesByUserID(t *testing.T) {
	cases := []struct {
		name           string
		path           string
		withAuth       bool
		setup          func(t *testing.T, suite *acceptanceSuite)
		svc            port.DeviceProfileService
		expectedStatus int
		assertFn       func(t *testing.T, status int, payload []byte, suite *acceptanceSuite)
	}{
		{
			name:     "returns profiles ordered by created_at desc",
			path:     "/device-profiles",
			withAuth: true,
			setup: func(t *testing.T, suite *acceptanceSuite) {
				for i := 0; i < 3; i++ {
					dp := entity.DeviceProfile{
						UserID:     suite.userID,
						Name:       fmt.Sprintf("profile_%d", i),
						DeviceType: "desktop",
					}
					require.NoError(t, suite.repo.CreateDeviceProfile(&dp))
					time.Sleep(5 * time.Millisecond)
				}
			},
			expectedStatus: nethttp.StatusOK,
			assertFn: func(t *testing.T, status int, payload []byte, suite *acceptanceSuite) {
				require.Equal(t, nethttp.StatusOK, status)
				var out []httpadapter.DeviceProfileResponse
				require.NoError(t, json.Unmarshal(payload, &out))
				assert.Len(t, out, 3)
				assert.Equal(t, "profile_2", out[0].Name)
				assert.Equal(t, suite.userID, out[0].UserID)
			},
		},
		{
			name:     "supports pagination",
			path:     "/device-profiles?page=2&page_size=1",
			withAuth: true,
			setup: func(t *testing.T, suite *acceptanceSuite) {
				for i := 0; i < 3; i++ {
					dp := entity.DeviceProfile{
						UserID:     suite.userID,
						Name:       fmt.Sprintf("page_profile_%d", i),
						DeviceType: "desktop",
					}
					require.NoError(t, suite.repo.CreateDeviceProfile(&dp))
					time.Sleep(5 * time.Millisecond)
				}
			},
			expectedStatus: nethttp.StatusOK,
			assertFn: func(t *testing.T, status int, payload []byte, _ *acceptanceSuite) {
				require.Equal(t, nethttp.StatusOK, status)
				var out []httpadapter.DeviceProfileResponse
				require.NoError(t, json.Unmarshal(payload, &out))
				require.Len(t, out, 1)
				assert.Equal(t, "page_profile_1", out[0].Name)
			},
		},
		{
			name:           "returns empty slice when user has no profiles",
			path:           "/device-profiles",
			withAuth:       true,
			expectedStatus: nethttp.StatusOK,
			assertFn: func(t *testing.T, status int, payload []byte, _ *acceptanceSuite) {
				require.Equal(t, nethttp.StatusOK, status)
				var out []httpadapter.DeviceProfileResponse
				require.NoError(t, json.Unmarshal(payload, &out))
				assert.Len(t, out, 0)
			},
		},
		{
			name:           "invalid page yields bad request",
			path:           "/device-profiles?page=0",
			withAuth:       true,
			expectedStatus: nethttp.StatusBadRequest,
			assertFn: func(t *testing.T, status int, payload []byte, _ *acceptanceSuite) {
				require.Equal(t, nethttp.StatusBadRequest, status)
				var out map[string]string
				require.NoError(t, json.Unmarshal(payload, &out))
				assert.Equal(t, "INVALID_ARGUMENT", out["code"])
			},
		},
		{
			name:           "invalid page size yields bad request",
			path:           "/device-profiles?page_size=0",
			withAuth:       true,
			expectedStatus: nethttp.StatusBadRequest,
			assertFn: func(t *testing.T, status int, payload []byte, _ *acceptanceSuite) {
				require.Equal(t, nethttp.StatusBadRequest, status)
				var out map[string]string
				require.NoError(t, json.Unmarshal(payload, &out))
				assert.Equal(t, "INVALID_ARGUMENT", out["code"])
			},
		},
		{
			name:     "maps domain errors to http codes",
			path:     "/device-profiles",
			withAuth: true,
			svc: &erroringDeviceProfileService{
				err: apperr.NewInvalidArgErr("bad", nil),
			},
			expectedStatus: nethttp.StatusBadRequest,
			assertFn: func(t *testing.T, status int, payload []byte, _ *acceptanceSuite) {
				require.Equal(t, nethttp.StatusBadRequest, status)
				var out map[string]string
				require.NoError(t, json.Unmarshal(payload, &out))
				assert.Equal(t, "INVALID_ARGUMENT", out["code"])
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			suite := newDeviceProfileSuite(t, tc.withAuth, tc.svc)
			if tc.setup != nil {
				tc.setup(t, suite)
			}

			headers := map[string]string{}
			if tc.withAuth {
				headers["Authorization"] = basicAuthHeader
			}

			resp := suite.doGet(t, tc.path, headers)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, resp.StatusCode)
			tc.assertFn(t, resp.StatusCode, body, suite)
		})
	}
}

func TestCreateDeviceProfile(t *testing.T) {
	type setupFn func(t *testing.T, suite *acceptanceSuite) ([]byte, map[string]string)
	cases := []struct {
		name           string
		setup          setupFn
		expectedStatus int
		assertFn       func(t *testing.T, status int, payload []byte, suite *acceptanceSuite)
	}{
		{
			name: "creates profile",
			setup: func(t *testing.T, suite *acceptanceSuite) ([]byte, map[string]string) {
				payload := map[string]any{
					"name":        "My profile",
					"device_type": "desktop",
				}
				body, err := json.Marshal(payload)
				require.NoError(t, err)
				return body, map[string]string{"Authorization": basicAuthHeader}
			},
			expectedStatus: nethttp.StatusCreated,
			assertFn: func(t *testing.T, status int, payload []byte, suite *acceptanceSuite) {
				require.Equal(t, nethttp.StatusCreated, status)
				var out httpadapter.DeviceProfileResponse
				require.NoError(t, json.Unmarshal(payload, &out))
				assert.Equal(t, "My profile", out.Name)
				assert.Equal(t, suite.userID, out.UserID)

				listed, err := suite.repo.ListDeviceProfiles(suite.userID.String(), 1, 10)
				require.NoError(t, err)
				assert.Equal(t, 1, len(listed))
			},
		},
		{
			name: "invalid json payload",
			setup: func(t *testing.T, suite *acceptanceSuite) ([]byte, map[string]string) {
				return []byte("{"), map[string]string{"Authorization": basicAuthHeader}
			},
			expectedStatus: nethttp.StatusBadRequest,
			assertFn: func(t *testing.T, status int, _ []byte, _ *acceptanceSuite) {
				require.Equal(t, nethttp.StatusBadRequest, status)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			suite := newDeviceProfileSuite(t, true, nil)
			body, headers := tc.setup(t, suite)
			resp := suite.doPost(t, "/device-profiles", headers, body)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, resp.StatusCode)
			tc.assertFn(t, resp.StatusCode, respBody, suite)
		})
	}
}

func TestUpdateDeviceProfile(t *testing.T) {
	cases := []struct {
		name           string
		setup          func(t *testing.T, suite *acceptanceSuite) (string, []byte)
		expectedStatus int
		assertFn       func(t *testing.T, status int, payload []byte)
	}{
		{
			name: "updates profile",
			setup: func(t *testing.T, suite *acceptanceSuite) (string, []byte) {
				dp := entity.DeviceProfile{UserID: suite.userID, Name: "Original", DeviceType: "desktop"}
				require.NoError(t, suite.repo.CreateDeviceProfile(&dp))
				payload := map[string]any{"name": "Updated"}
				body, err := json.Marshal(payload)
				require.NoError(t, err)
				return dp.ID.String(), body
			},
			expectedStatus: nethttp.StatusOK,
			assertFn: func(t *testing.T, status int, payload []byte) {
				require.Equal(t, nethttp.StatusOK, status)
				var out httpadapter.DeviceProfileResponse
				require.NoError(t, json.Unmarshal(payload, &out))
				assert.Equal(t, "Updated", out.Name)
			},
		},
		{
			name: "rejects empty payload",
			setup: func(t *testing.T, suite *acceptanceSuite) (string, []byte) {
				dp := entity.DeviceProfile{UserID: suite.userID, Name: "Original", DeviceType: "desktop"}
				require.NoError(t, suite.repo.CreateDeviceProfile(&dp))
				return dp.ID.String(), []byte("{}")
			},
			expectedStatus: nethttp.StatusBadRequest,
			assertFn: func(t *testing.T, status int, _ []byte) {
				require.Equal(t, nethttp.StatusBadRequest, status)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			suite := newDeviceProfileSuite(t, true, nil)
			id, body := tc.setup(t, suite)

			headers := map[string]string{"Authorization": basicAuthHeader}
			resp := suite.doPut(t, fmt.Sprintf("/device-profiles/%s", id), headers, body)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, resp.StatusCode)
			tc.assertFn(t, resp.StatusCode, respBody)
		})
	}
}

func TestDeleteDeviceProfile(t *testing.T) {
	cases := []struct {
		name           string
		setup          func(t *testing.T, suite *acceptanceSuite) string
		expectedStatus int
		assertFn       func(t *testing.T, status int, suite *acceptanceSuite)
	}{
		{
			name: "deletes profile",
			setup: func(t *testing.T, suite *acceptanceSuite) string {
				dp := entity.DeviceProfile{UserID: suite.userID, Name: "ToDelete", DeviceType: "desktop"}
				require.NoError(t, suite.repo.CreateDeviceProfile(&dp))
				return dp.ID.String()
			},
			expectedStatus: nethttp.StatusNoContent,
			assertFn: func(t *testing.T, status int, suite *acceptanceSuite) {
				require.Equal(t, nethttp.StatusNoContent, status)
				remaining, err := suite.repo.ListDeviceProfiles(suite.userID.String(), 1, 10)
				require.NoError(t, err)
				assert.Len(t, remaining, 0)
			},
		},
		{
			name: "invalid id returns bad request",
			setup: func(t *testing.T, suite *acceptanceSuite) string {
				return "bad-id"
			},
			expectedStatus: nethttp.StatusBadRequest,
			assertFn: func(t *testing.T, status int, _ *acceptanceSuite) {
				require.Equal(t, nethttp.StatusBadRequest, status)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			suite := newDeviceProfileSuite(t, true, nil)
			id := tc.setup(t, suite)

			headers := map[string]string{"Authorization": basicAuthHeader}
			resp := suite.doDelete(t, fmt.Sprintf("/device-profiles/%s", id), headers)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatus, resp.StatusCode)
			tc.assertFn(t, resp.StatusCode, suite)
		})
	}
}

type erroringDeviceProfileService struct {
	err error
}

func (e *erroringDeviceProfileService) ListDeviceProfilesByUserID(context.Context, int, int) ([]entity.DeviceProfile, error) {
	return nil, e.err
}

func (e *erroringDeviceProfileService) CreateDeviceProfile(context.Context, *entity.DeviceProfile) error {
	return fmt.Errorf("not implemented")
}

func (e *erroringDeviceProfileService) UpdateDeviceProfile(context.Context, *entity.DeviceProfile) (*entity.DeviceProfile, error) {
	return nil, fmt.Errorf("not implemented")
}

func (e *erroringDeviceProfileService) DeleteDeviceProfile(context.Context, string) error {
	return fmt.Errorf("not implemented")
}
