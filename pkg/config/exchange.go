package config

import (
	"context"
	"fmt"

	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/connector/binance"
	"github.com/nbitslabs/nOracle/pkg/connector/bybit"
)

func LoadConnector(ctx context.Context, exchange string, symbols []string) (connector.ExchangeConnector, error) {
	switch exchange {
	case "binance":
		return binance.NewConnector(ctx, symbols)
	case "bybit":
		return bybit.NewConnector(ctx, symbols)
	default:
		return nil, fmt.Errorf("unknown exchange: %s", exchange)
	}
}
