package repo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Nadzyoki/BannerCounter/internal/utils"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) SaveCount(ctx context.Context, ids []string, counts []int) error {
	if len(ids) != len(counts) {
		return errors.New(ErrLenMismatch)
	}
	if len(ids) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrBeginTx, err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		}
	}()

	_, err = tx.Exec(ctx, `
		CREATE TEMP TABLE tmp_banner_stats (
			ts_minute TIMESTAMP WITHOUT TIME ZONE,
			banner_id BIGINT,
			count BIGINT
		) ON COMMIT DROP
	`)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", ErrCreateTempTable, err)
	}

	rows := make([][]interface{}, 0, len(ids))
	for i, key := range ids {
		bannerIDStr, ts, err := utils.SplitKey(key)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("%s: %w", ErrSplitKey, err)
		}

		bid, err := strconv.ParseInt(bannerIDStr, 10, 64)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("%s: %w", ErrParseBannerID, err)
		}

		rows = append(rows, []interface{}{ts, bid, counts[i]})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"tmp_banner_stats"},
		[]string{"ts_minute", "banner_id", "count"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", ErrCopyFrom, err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO banner_stats (ts_minute, banner_id, count)
		SELECT ts_minute, banner_id, count FROM tmp_banner_stats
		ON CONFLICT (ts_minute, banner_id)
		DO UPDATE SET count = banner_stats.count + EXCLUDED.count
	`)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("%s: %w", ErrMergeIntoMainTable, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", ErrCommit, err)
	}

	return nil
}
