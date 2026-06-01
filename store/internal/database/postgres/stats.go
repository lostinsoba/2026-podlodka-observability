package postgres

import (
	"context"

	"store/internal/model"
)

func (d *Database) EvaluateMessagesStats(ctx context.Context) (model.MessagesStats, error) {
	const (
		query = `
			select n_live_tup
			from pg_stat_user_tables
			where relname = 'store_message'`
	)
	var count int64
	err := d.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return model.MessagesStats{}, err
	}
	return model.MessagesStats{MessageCount: count}, nil
}
