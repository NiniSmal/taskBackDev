package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"os"
	"taskBackDev/api"
	"taskBackDev/config"
	"taskBackDev/service"
	"taskBackDev/storage"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error("get config", "err", err)
	}

	conn, err := pgxpool.New(ctx, cfg.Postgres)
	if err != nil {
		logger.Error("connect postgres", "err", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		logger.Error("unable to ping to database", "error", err)
	}
	logger.Info("Connected to Postgres OK")

	defer conn.Close()
	st := storage.NewStorage(conn)
	server := service.NewService(st, cfg)
	handler := api.NewHandler(server)
	r := gin.Default()
	r.GET("/sing_in", handler.SingIn)
	r.GET("/refresh", handler.Refresh)
	err = r.Run(fmt.Sprintf("localhost:%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
}
