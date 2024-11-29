package exchange_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type repo struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func New(db *pgxpool.Pool, logger *zap.Logger) *repo {
	return &repo{
		db:     db,
		logger: logger,
	}
}

func (r *repo) WriteExchange(ctx context.Context, askPrice float64, bidPrice float64, createdAt string) error {
	timestamp, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		r.logger.Error("failed to parse timestamp", zap.Error(err))
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	query := `
        INSERT INTO rates (market, ask_price, bid_price, timestamp)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (market, timestamp) DO NOTHING
    `

	_, err = r.db.Exec(ctx, query, "usdtrub", askPrice, bidPrice, timestamp)
	if err != nil {
		r.logger.Error("failed to insert exchange rate", zap.Error(err))
		return fmt.Errorf("failed to insert exchange rate: %w", err)
	}

	r.logger.Info("successfully inserted exchange rate", zap.String("market", "usdtrub"))
	return nil
}
