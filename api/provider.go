package main

import (
	"context"
	"fmt"

	"github.com/0x726f6f6b6965/task/internal/config"
	zaplog "github.com/0x726f6f6b6965/task/internal/log"
	"github.com/0x726f6f6b6965/task/internal/services"
	"github.com/0x726f6f6b6965/task/internal/utils"
	pbTask "github.com/0x726f6f6b6965/task/protos/task/v1"
	"github.com/google/wire"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

var applicationSet = wire.NewSet(componentSet, services.NewTaskService, newServer)

var componentSet = wire.NewSet(generatorSet, loggerSet, dbSet)

var loggerSet = wire.NewSet(logCfg, zaplog.NewLogger)

var dbSet = wire.NewSet(redisCfg, redisClient)

var generatorSet = wire.NewSet(generatorCfg, utils.NewGenerator)

func logCfg(cfg *config.Config) *config.Log {
	return &cfg.Log
}

func redisCfg(cfg *config.Config) *redis.Options {
	return &redis.Options{
		Addr:       fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Username:   cfg.Redis.User,
		Password:   cfg.Redis.Password,
		DB:         cfg.Redis.DB,
		MaxRetries: cfg.Redis.MaxRetries,
	}
}

func generatorCfg(cfg *config.Config) uint64 {
	return cfg.NodeID
}

func redisClient(opt *redis.Options) (*redis.Client, func(), error) {
	client := redis.NewClient(opt)
	return client, func() { client.Close() }, nil
}

func newServer(ctx context.Context, server pbTask.TaskServiceServer, logger *zap.Logger) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseEnumNumbers:  true,
		},
	}))

	err := pbTask.RegisterTaskServiceHandlerServer(ctx, mux, server)
	if err != nil {
		logger.Error("failed to register", zap.Error(err))
		return mux, err
	}
	return mux, nil
}
