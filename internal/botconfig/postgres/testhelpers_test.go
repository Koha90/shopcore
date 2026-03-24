package postgres

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/koha90/shopcore/pkg/migrator"
)

func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	container, err := tcpostgres.Run(
		ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("github.com/koha90/shopcore_test"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
	)
	if err != nil {
		t.Fatalf("run postgres container: %v", err)
	}

	t.Cleanup(func() {
		_ = container.Terminate(context.Background())
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("postgres connection string: %v", err)
	}

	if err := waitForPostgres(ctx, dsn, 10*time.Second); err != nil {
		t.Fatalf("wait for postgres: %v", err)
	}

	migrationsPath := resolveMigrationsPath(t)

	if err = migrator.MigratePostgres(dsn, migrationsPath); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("create pgx pool: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		t.Fatalf("ping pgx pool: %v", err)
	}

	t.Cleanup(pool.Close)

	return pool
}

// dbExecutor is a minimal test-only abstraction used by fixture helpers.
// It keep helper code small and decoupled from concrete pool type.
// type dbExecutor interface {
// 	Exec(ctx context.Context, sql string, arguments ...any) (commandTag, error)
// }

func waitForPostgres(ctx context.Context, dsn string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for {
		pool, err := pgxpool.New(ctx, dsn)
		if err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, time.Second)
			pingErr := pool.Ping(pingCtx)
			cancel()
			pool.Close()

			if pingErr == nil {
				return nil
			}
			err = pingErr
		}

		if time.Now().After(deadline) {
			return err
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func resolveMigrationsPath(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("cannot resolve current file path")
	}

	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", ".."))
	return filepath.Join(projectRoot, "migrations")
}
