package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	middleswares "github.com/nahtann/controlriver.com/api/v1/middlewares"
	routes_v1 "github.com/nahtann/controlriver.com/api/v1/routes"
	"github.com/nahtann/controlriver.com/internal/helpers/migrations"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	migrations.Up()

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	repository := repository.New(dbpool)

	logger := NewLogger()
	defer logger.Sync()

	middlewareOrchestrator := middleswares.NewMiddlewareOrchestration(logger)

	routerV1 := routes_v1.NewRouterV1(middlewareOrchestrator)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", routerV1))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server on ", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

	fmt.Println("Server stopped")
}

func NewLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stderr),
		zap.InfoLevel,
	)

	return zap.New(core)
}
