package connector

import (
	"context"
	"math/big"
	"strings"
)

type ExchangeConnector interface {
	Name() Exchange
	Tickers() []Symbol
	Close() error
	StreamTickers(ctx context.Context, out chan<- TickerUpdate) error
}

type TickerUpdate struct {
	Exchange  Exchange
	Symbol    Symbol
	Price     *big.Float
	Volume    *big.Float
	Timestamp int64
}

type Exchange string
type Symbol string

func (s Symbol) String() string {
	return strings.ToLower(string(s))
}
