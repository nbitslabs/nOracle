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
	Exchange  string     `json:"exchange"`
	Symbol    string     `json:"symbol"`
	Price     *big.Float `json:"price"`
	Volume    *big.Float `json:"volume"`
	Timestamp int64      `json:"timestamp"`
}
