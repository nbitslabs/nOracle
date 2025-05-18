package oracle

import (
	"context"
	"log/slog"

	"github.com/nbitslabs/nOracle/pkg/config"
	"github.com/nbitslabs/nOracle/pkg/connector"
)

type Services struct {
	Exchanges []connector.ExchangeConnector
}

func NewServices(ctx context.Context, configPath string) (*Services, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	exchanges := []connector.ExchangeConnector{}
	for _, exchange := range c.Exchanges {
		slog.Info("Loading exchange", "name", exchange.Name, "url", exchange.URL, "symbols", exchange.Symbols)
		conn, err := config.LoadConnector(ctx, exchange.Name, exchange.URL, exchange.Symbols)
		if err != nil {
			return nil, err
		}

		exchanges = append(exchanges, conn)
	}

	return &Services{
		Exchanges: exchanges,
	}, nil
}
