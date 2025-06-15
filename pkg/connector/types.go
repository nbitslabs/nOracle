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
	Exchange  string `json:"exchange"`
	Symbol    string `json:"symbol"`
	Timestamp int64  `json:"timestamp"`

	Spot    *SpotPriceUpdate    `json:"spot,omitempty"`
	Futures *FuturesPriceUpdate `json:"futures,omitempty"`
}

type SpotPriceUpdate struct {
	Price  *big.Float `json:"price"`
	Volume *big.Float `json:"volume"`
}

type FuturesPriceUpdate struct {
	MarkPrice   *big.Float `json:"markPrice"`
	IndexPrice  *big.Float `json:"indexPrice"`
	FundingRate *big.Float `json:"fundingRate"`
}
