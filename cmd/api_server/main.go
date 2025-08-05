package main

import (
	"fmt"
	"net/http"
	"os"

	middleswares "github.com/nahtann/controlriver.com/api/v1/middlewares"
	routes_v1 "github.com/nahtann/controlriver.com/api/v1/routes"
	"github.com/nahtann/controlriver.com/internal/helpers/migrations"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	migrations.Up()

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
