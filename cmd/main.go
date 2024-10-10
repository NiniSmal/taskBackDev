package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"net/http"
	"os"
	"taskBackDev/api"
	"taskBackDev/config"
	"taskBackDev/service"
	"taskBackDev/service/send_service"
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
	sendService := send_service.NewSendService(cfg.Email, cfg.Password)
	handler := api.NewHandler(server, cfg, sendService)
	r := http.NewServeMux()

	r.HandleFunc("GET /auth", handler.AuthHandler())
	r.HandleFunc("GET /refresh", handler.Refresh)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
	if err != nil {
		log.Fatal(err)
	}
}
