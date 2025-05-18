package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/nbitslabs/nOracle/internal/oracle"
	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/utils/env"
)

var configPath = env.Get("CONFIG_PATH", "config.yaml")

func main() {
	slog.Info("Starting nOracle", "config_path", configPath)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	exit := make(chan struct{})

	service, err := oracle.NewServices(context.Background(), configPath)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	out := make(chan connector.TickerUpdate)
	for _, exchange := range service.Exchanges {
		slog.Info("Streaming tickers", "name", exchange.Name())
		go exchange.StreamTickers(context.Background(), out)
		defer exchange.Close()
	}

	slog.Info("nOracle started")
	<-sig
	close(exit)
}
