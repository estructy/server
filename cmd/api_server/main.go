package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	middleswares "github.com/nahtann/controlriver.com/api/v1/middlewares"
	routesv1 "github.com/nahtann/controlriver.com/api/v1/routes"
	"github.com/nahtann/controlriver.com/internal/helpers/migrations"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @todo: refactor this
type myQueryTracer struct {
	log *zap.Logger
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	tracer.log.Info("Executing command",
		zap.String("sql", data.SQL),
		zap.Any("args", data.Args),
	)

	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func main() {
	// --- Logger setup ---
	logger := NewLogger()
	defer logger.Sync()
	// --------

	// --- Database setup ---
	migrations.Up()
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}
	config.ConnConfig.Tracer = &myQueryTracer{
		log: logger.WithOptions(zap.AddCallerSkip(1)),
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()
	repository := repository.New(dbpool)
	// ------

	// --- Router setup ---
	middlewareOrchestrator := middleswares.NewMiddlewareOrchestration(logger)
	routerV1 := routesv1.NewRouterV1(middlewareOrchestrator, dbpool, repository)
	// ------

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", routerV1))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Starting server on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Error starting server: %v\n", err)
		return
	}

	log.Printf("Server stopped\n")
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
