package config

import (
	"context"
	"fmt"

	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/connector/binance"
	"github.com/nbitslabs/nOracle/pkg/connector/okx"
)

func LoadConnector(ctx context.Context, exchange string, url string, symbols []string) (connector.ExchangeConnector, error) {
	switch exchange {
	case "binance":
		return binance.NewConnector(ctx, url, symbols)
	case "okx":
		return okx.NewConnector(ctx, url, symbols)
	default:
		return nil, fmt.Errorf("unknown exchange: %s", exchange)
	}
}
