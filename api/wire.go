//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/0x726f6f6b6965/task/internal/config"
	"github.com/google/wire"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func initApplication(ctx context.Context, cfg *config.Config) (*runtime.ServeMux, func(), error) {
	panic(wire.Build(applicationSet))
}
