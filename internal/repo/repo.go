package repo

import (
	"context"
	"fmt"

	"github.com/Nadzyoki/BannerCounter/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ErrLenMismatch        = "len(ids) != len(counts)"
	ErrBeginTx            = "begin tx error"
	ErrCreateTempTable    = "create temp table error"
	ErrSplitKey           = "split key error"
	ErrParseBannerID      = "parse bannerID error"
	ErrCopyFrom           = "copy from error"
	ErrMergeIntoMainTable = "merge into main table error"
	ErrCommit             = "commit error"
	ErrMigration          = "migration error"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(ctx context.Context, cfg *config.Config) (*Repo, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DB,
		cfg.Password,
	)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	err = runMigrations(cfg)
	if err != nil {
		return nil, err
	}

	return &Repo{db: pool}, nil
}

func (r *Repo) Close() {
	r.db.Close()
}
