package oracle

import (
	"context"
	"log/slog"

	"github.com/nbitslabs/nOracle/pkg/config"
	"github.com/nbitslabs/nOracle/pkg/connector"
)

type Services struct {
	Exchanges []connector.ExchangeConnector
	Info      *Info
}

type Info struct {
	AvailableExchanges []string
	AvailableSymbols   map[string][]string
}

func NewServices(ctx context.Context, configPath string) (*Services, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	exchanges := []connector.ExchangeConnector{}
	availableExchanges := []string{}
	availableSymbols := map[string][]string{}
	for _, exchange := range c.Exchanges {
		slog.Info("Loading exchange", "name", exchange.Name, "symbols", exchange.Symbols)
		conn, err := config.LoadConnector(ctx, exchange.Name, exchange.Symbols)
		if err != nil {
			return nil, err
		}

		exchanges = append(exchanges, conn)
		availableExchanges = append(availableExchanges, exchange.Name)
		for _, symbol := range exchange.Symbols {
			availableSymbols[exchange.Name] = append(availableSymbols[exchange.Name], symbol)
		}
	}

	return &Services{
		Exchanges: exchanges,
		Info: &Info{
			AvailableExchanges: availableExchanges,
			AvailableSymbols:   availableSymbols,
		},
	}, nil
}
