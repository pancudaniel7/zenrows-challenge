package util

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/viper"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitTestContainers starts a Postgres 16 container for integration tests,
// aligns credentials with configs/test.yml (or env overrides), updates viper
// database settings to point to the container, and applies the SQL init scripts
// from deployments/postgres. It registers a cleanup with t.Cleanup.
func InitTestContainers(t *testing.T) (func(), error) {
	t.Helper()

	// Resolve desired DB credentials
	dbName, dbUser, dbPass := desiredDBConfig()

    // Start postgres container and bind to the configured port (from test.yml)
    ctx := context.Background()
    basePort := viper.GetInt("database.port")
    desiredHostPort := 0
    if basePort > 0 { desiredHostPort = basePort }
    c, host, port, err := startPostgres16(ctx, dbName, dbUser, dbPass, desiredHostPort)
	if err != nil {
		return nil, err
	}

	// Point config to container
	pointViperToContainer(host, port, dbName, dbUser, dbPass)

	// Apply schema and seed scripts
	if err := applyInitScripts(host, port, dbName, dbUser, dbPass); err != nil {
		_ = c.Terminate(ctx)
		return nil, err
	}

	cleanup := func() { _ = c.Terminate(ctx) }
	t.Cleanup(cleanup)
	return cleanup, nil
}

// desiredDBConfig reads DB name/user/pass from viper with sensible defaults.
func desiredDBConfig() (name, user, pass string) {
	name = viper.GetString("database.database")
	user = viper.GetString("database.user")
	pass = viper.GetString("database.password")
	return
}

// startPostgres16 runs a Postgres 16 container and waits for readiness.
func startPostgres16(ctx context.Context, dbName, dbUser, dbPass string, hostPort int) (tc.Container, string, int, error) {
	req := tc.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPass,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(1),
			wait.ForListeningPort("5432/tcp"),
		).WithDeadline(2 * time.Minute),
	}
	if hostPort > 0 {
		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				nat.Port("5432/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(hostPort)}},
			}
		}
	}
	c, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{ContainerRequest: req, Started: true})
	if err != nil {
		return nil, "", 0, err
	}
	host, err := c.Host(ctx)
	if err != nil {
		_ = c.Terminate(ctx)
		return nil, "", 0, err
	}
	mp, err := c.MappedPort(ctx, "5432/tcp")
	if err != nil {
		_ = c.Terminate(ctx)
		return nil, "", 0, err
	}
	return c, host, mp.Int(), nil
}

// pointViperToContainer sets viper's database.* keys to provided host/port/creds.
func pointViperToContainer(host string, port int, dbName, dbUser, dbPass string) {
	viper.Set("database.host", host)
	viper.Set("database.port", port)
	viper.Set("database.database", dbName)
	viper.Set("database.user", dbUser)
	viper.Set("database.password", dbPass)
	viper.Set("database.sslmode", "disable")
}

// applyInitScripts connects to the new DB and executes all SQL files from
// deployments/postgres in lexical order.
func applyInitScripts(host string, port int, dbName, dbUser, dbPass string) error {
    dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbName, dbUser, dbPass)
    gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }
    scriptsDir, err := locateScriptsDir()
    if err != nil {
        return err
    }
    entries, err := os.ReadDir(scriptsDir)
    if err != nil {
        return err
    }
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) == ".sql" {
			files = append(files, filepath.Join(scriptsDir, name))
		}
	}
	sort.Strings(files)
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if err := gdb.Exec(string(b)).Error; err != nil {
			return fmt.Errorf("apply %s: %w", filepath.Base(f), err)
		}
	}
    return nil
}

// locateScriptsDir finds the deployments/postgres directory robustly by checking
// multiple candidate paths relative to config file location and working dir.
func locateScriptsDir() (string, error) {
    var candidates []string
    if cfg := viper.ConfigFileUsed(); cfg != "" {
        base := filepath.Dir(cfg)
        candidates = append(candidates,
            filepath.Join(base, "deployments", "postgres"),
            filepath.Join(base, "..", "deployments", "postgres"),
            filepath.Join(base, "..", "..", "deployments", "postgres"),
        )
    }
    if wd, err := os.Getwd(); err == nil {
        candidates = append(candidates,
            filepath.Join(wd, "deployments", "postgres"),
            filepath.Join(wd, "..", "deployments", "postgres"),
            filepath.Join(wd, "..", "..", "deployments", "postgres"),
        )
    }
    // Fallback
    candidates = append(candidates, filepath.Join("deployments", "postgres"))

    for _, p := range candidates {
        if fi, err := os.Stat(p); err == nil && fi.IsDir() {
            return p, nil
        }
    }
    return "", fmt.Errorf("deployments/postgres directory not found")
}
