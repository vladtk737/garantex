package main

import (
	"context"
	"fmt"
	"garantex/config"
	grpcServer "garantex/internal/application/grpc"
	"garantex/internal/domain/usecase"
	"garantex/internal/infrastructure/http"
	repo "garantex/internal/infrastructure/repository/exchange_repo"
	"garantex/pkg/logger"
	"garantex/pkg/psql"
	grpcSchema "garantex/schemas/proto/grpc/garantex/schemas/proto/grpc"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
)

const (
	exitStatusOk     = 0
	exitStatusFailed = 1
)

func main() {
	// загружает переменные окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config init error: %s", err)
	}

	os.Exit(run(cfg))
}

func run(cfg *config.Config) (exitStatus int) {
	logConfig := logger.LogConfig{
		Level:       cfg.Log.Level,
		Encoding:    cfg.Log.Encoding,
		OutputPaths: cfg.Log.OutputPaths,
		ErrorOutput: cfg.Log.ErrorOutput,
	}

	logger, err := logger.NewZapLogger(logConfig)
	if err != nil {
		log.Fatal("logger init error", zap.Error(err))
	}

	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Error("recover after panic", zap.Any("err", panicErr), zap.String("stacktrace", string(debug.Stack())))
			exitStatus = exitStatusFailed
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	db, err := psql.Connect(ctx,
		psql.WithHostPort(cfg.PG.Host),
		psql.WithPort(cfg.PG.Port),
		psql.WithDatabase(cfg.PG.Database),
		psql.WithUser(cfg.PG.User),
		psql.WithPassword(cfg.PG.Password),
		psql.WithUserAdmin(cfg.PG.UserAdmin),
		psql.WithPasswordAdmin(cfg.PG.PasswordAdmin),
		psql.WithMigrations(os.DirFS("db/migrations")),
		psql.WithLogger(logger),
	)
	if err != nil {
		logger.Error("db init error", zap.Error(err))
		return exitStatusFailed
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		logger.Error("grpc failed to listen", zap.Error(err))
		return exitStatusFailed
	}

	httpClient := http.New("https://garantex.org")
	exchangeRepo := repo.New(db, logger)
	processTradesUsecase := usecase.New(exchangeRepo, httpClient, logger)
	serviceHandler := grpcServer.New(processTradesUsecase, logger)

	healthServer := &grpcServer.HealthServer{Db: db}
	grpcSrv := grpc.NewServer()

	grpcSchema.RegisterExchangeServiceServer(grpcSrv, serviceHandler)
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthServer)

	errChan := make(chan error)

	go func() {
		logger.Info("Starting grpc server...")
		if err := grpcSrv.Serve(lis); err != nil {
			logger.Error("grpc server init error", zap.Error(err))
			errChan <- err
		}
	}()

	defer func() {
		grpcSrv.GracefulStop()
	}()

	select {
	case err := <-errChan:
		logger.Error("fatal error, service shutdown", zap.Error(err))
	case <-ctx.Done():
		logger.Info("service shutdown")
	}

	return exitStatusOk
}
