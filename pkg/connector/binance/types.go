package binance

import (
	"context"
	"math/big"

	"github.com/recws-org/recws"
)

type Connector struct {
	ctx   context.Context
	pairs []string

	ws *recws.RecConn
}

const Name = "binance"

type Ticker struct {
	Name   string       `json:"stream"`
	Stream TickerStream `json:"data"`
}

type TickerStream struct {
	EventType          string     `json:"e"`
	EventTime          int64      `json:"E"`
	Symbol             string     `json:"s"`
	PriceChange        *big.Float `json:"p"`
	PriceChangePercent *big.Float `json:"P"`
	WeightedAvgPrice   *big.Float `json:"w"`
	FirstTradePrice    *big.Float `json:"x"`
	LastPrice          *big.Float `json:"c"`
	LastQuantity       *big.Float `json:"Q"`
	BestBidPrice       *big.Float `json:"b"`
	BestBidQuantity    *big.Float `json:"B"`
	BestAskPrice       *big.Float `json:"a"`
	BestAskQuantity    *big.Float `json:"A"`
	OpenPrice          *big.Float `json:"o"`
	HighPrice          *big.Float `json:"h"`
	LowPrice           *big.Float `json:"l"`
	Volume             *big.Float `json:"v"`
	QuoteVolume        *big.Float `json:"q"`
	OpenTime           *big.Int   `json:"O"`
	CloseTime          *big.Int   `json:"C"`
	FirstTradeID       *big.Int   `json:"F"`
	LastTradeID        *big.Int   `json:"L"`
	TotalTrades        *big.Int   `json:"n"`
}
