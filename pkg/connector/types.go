package connector

import (
	"context"
	"math/big"
)

type ExchangeConnector interface {
	Name() string
	Tickers() []string
	Close() error
	StreamTickers(ctx context.Context, out chan<- TickerUpdate) error
}

type TickerUpdate struct {
	Exchange  string
	Symbol    string
	Price     *big.Float
	Volume    *big.Float
	Timestamp int64
}
