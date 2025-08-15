package repo

import (
	"context"
	"time"
)

func (r *Repo) GetStats(ctx context.Context, bannerID string, tsFrom, tsTo string) (map[string]int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT ts_minute, count 
		FROM banner_stats
		WHERE banner_id = $1
		AND ts_minute >= $2
		AND ts_minute <= $3
		ORDER BY ts_minute
	`, bannerID, tsFrom, tsTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var ts time.Time
		var count int
		if err := rows.Scan(&ts, &count); err != nil {
			return nil, err
		}
		result[ts.Format(time.RFC3339)] = count
	}
	return result, nil
}
