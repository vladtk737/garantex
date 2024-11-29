package usecase

import (
	"context"
	"fmt"
	"garantex/internal/infrastructure/http"
	proto "garantex/schemas/proto/grpc/garantex/schemas/proto/grpc"
	"go.uber.org/zap"
	"strconv"
)

type repo interface {
	WriteExchange(ctx context.Context, askPrice float64, bidPrice float64, createdAt string) error
}

type Usecase struct {
	repo     repo
	garantex *http.GarantexClient
	logger   *zap.Logger
}

func New(repo repo, garantex *http.GarantexClient, logger *zap.Logger) *Usecase {
	return &Usecase{
		repo:     repo,
		garantex: garantex,
		logger:   logger,
	}
}

func (u *Usecase) ProcessTrades(ctx context.Context) ([]*proto.TradePriceData, error) {
	u.logger.Info("fetching trades from Garantex API")

	trades, err := u.garantex.GetTrades()
	if err != nil {
		u.logger.Error("failed to fetch trades from API", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch trades from API: %w", err)
	}

	if len(trades) == 0 {
		u.logger.Warn("no trades received from API")
		return nil, fmt.Errorf("no trades received from API")
	}

	u.logger.Info("starting to process trades", zap.Int("number_of_trades", len(trades)))

	var maxPrice, minPrice float64
	var createdAt string

	maxPrice, err = strconv.ParseFloat(trades[0].Price, 64)
	if err != nil {
		u.logger.Error("failed to parse initial price", zap.Error(err))
		return nil, fmt.Errorf("failed to parse initial price: %w", err)
	}
	minPrice = maxPrice
	createdAt = trades[0].CreatedAt

	for _, trade := range trades {
		price, err := strconv.ParseFloat(trade.Price, 64)
		if err != nil {
			u.logger.Error("failed to parse price", zap.Int("trade_id", trade.ID), zap.Error(err))
			return nil, fmt.Errorf("failed to parse price: %w", err)
		}

		if price > maxPrice {
			maxPrice = price
		}
		if price < minPrice {
			minPrice = price
		}
	}

	err = u.repo.WriteExchange(ctx, maxPrice, minPrice, createdAt)
	if err != nil {
		u.logger.Error("failed to write exchange data to database", zap.Error(err))
		return nil, fmt.Errorf("failed to write exchange data to database: %w", err)
	}

	u.logger.Info("successfully processed trades", zap.Int("number_of_trades", len(trades)))

	var response []*proto.TradePriceData
	for _, trade := range trades {
		response = append(response, &proto.TradePriceData{
			Id:        int32(trade.ID),
			AskPrice:  fmt.Sprintf("%f", maxPrice),
			BidPrice:  fmt.Sprintf("%f", minPrice),
			Timestamp: createdAt,
		})
	}

	return response, nil
}
