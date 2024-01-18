package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/0x726f6f6b6965/task/internal/config"
	zaplog "github.com/0x726f6f6b6965/task/internal/log"
	"github.com/0x726f6f6b6965/task/internal/services"
	"github.com/0x726f6f6b6965/task/internal/utils"
	pbTask "github.com/0x726f6f6b6965/task/protos/task/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

func main() {
	godotenv.Load()
	path := os.Getenv("CONFIG")
	var cfg config.Config
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("read yaml error", err)
		return
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("unmarshal yaml error", err)
		return
	}

	zaplog, cleanup, err := zaplog.NewLogger(&cfg.Log)
	if err != nil {
		log.Fatal("create log error", err)
		return
	}
	defer cleanup()

	generator, err := utils.NewGenerator(cfg.NodeID)
	if err != nil {
		zaplog.Error("create generator fail", zap.Error(err))
		return
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Username:   cfg.Redis.User,
		Password:   cfg.Redis.Password,
		DB:         cfg.Redis.DB,
		MaxRetries: cfg.Redis.MaxRetries,
	})
	defer redisClient.Close()

	taskService := services.NewTaskService(generator, redisClient, zaplog)
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseEnumNumbers:  true,
		},
	}))
	err = pbTask.RegisterTaskServiceHandlerServer(context.Background(), mux, taskService)
	if err != nil {
		zaplog.Error("failed to register", zap.Error(err))
		return
	}

	zaplog.Info("server listening", zap.Int("port", cfg.Rest.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Rest.Port), mux); err != nil {
		zaplog.Error("failed to serve", zap.Error(err))
		return
	}
}
