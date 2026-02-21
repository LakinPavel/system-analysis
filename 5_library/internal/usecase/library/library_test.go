//nolint:errcheck //why not
package library

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/project/library/config"
	"github.com/project/library/internal/usecase/repository"
	"go.uber.org/zap"
)

func TestLibrary(t *testing.T) {
	t.Helper()
	t.Parallel()
	cfg, err := config.New()
	if err != nil {
		os.Exit(-1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	dbPool, err := pgxpool.New(ctx, cfg.PG.URL)

	if err != nil {
		return
	}

	defer dbPool.Close()
	forTests := New(zap.NewNop(), repository.NewPostgresRepository(dbPool), repository.NewPostgresRepository(dbPool))
	forTests.RegisterAuthor(context.Background(), "")
	forTests.GetAuthor(context.Background(), "")
	forTests.RegisterBook(context.Background(), "", []string{})
	forTests.GetBook(context.Background(), "")
	forTests.ChangeAuthorInfo(context.Background(), "", "")
	forTests.UpdateBook(context.Background(), "", "", []string{})
	forTests.GetAuthorBooks(context.Background(), "")
}
