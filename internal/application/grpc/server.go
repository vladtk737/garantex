package grpc

import (
	"context"
	"garantex/internal/domain/usecase"
	proto "garantex/schemas/proto/grpc/garantex/schemas/proto/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	proto.UnimplementedExchangeServiceServer
	Usecase *usecase.Usecase
	logger  *zap.Logger
}

func New(uc *usecase.Usecase, logger *zap.Logger) *Server {
	return &Server{
		Usecase: uc,
		logger:  logger,
	}
}

func (s *Server) GetTrades(ctx context.Context, empty *emptypb.Empty) (*proto.GetTradesResponse, error) {
	trades, err := s.Usecase.ProcessTrades(ctx)
	if err != nil {
		s.logger.Error("failed to process trades",
			zap.Error(err))

		return &proto.GetTradesResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, status.Errorf(codes.Internal, err.Error())
	}

	s.logger.Info("successfully processed trades")

	return &proto.GetTradesResponse{
		Success:        true,
		TradePriceData: trades,
	}, nil
}
