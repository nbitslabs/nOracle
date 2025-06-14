package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/nbitslabs/nOracle/internal/oracle"
	"github.com/nbitslabs/nOracle/internal/route/info"
	"github.com/nbitslabs/nOracle/internal/route/ticker"
	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/storage"
	"github.com/nbitslabs/nOracle/pkg/utils/env"
)

var configPath = env.Get("CONFIG_PATH", "config.yaml")

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
}

func main() {
	slog.Info("Starting nOracle", "config_path", configPath)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	service, err := oracle.NewServices(ctx, configPath)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	tickerStore := storage.NewMemory[connector.TickerUpdate]()
	out := make(chan connector.TickerUpdate)
	for _, exchange := range service.Exchanges {
		slog.Info("Streaming tickers", "name", exchange.Name())
		go exchange.StreamTickers(ctx, out)
		defer exchange.Close()
	}

	// Initialize APIs
	tickerAPI := ticker.NewAPI(ctx, tickerStore, out)
	infoAPI := info.NewAPI(ctx, service.Info)

	router := gin.Default()
	tickerAPI.Routes(router)
	infoAPI.Routes(router)

	slog.Info("nOracle started")
	go func() {
		if err := router.Run(":8080"); err != nil {
			slog.Error("Failed to start router", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
}
