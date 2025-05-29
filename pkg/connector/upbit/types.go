package upbit

import (
	"context"

	"github.com/recws-org/recws"
)

type Connector struct {
	ctx   context.Context
	pairs []string

	ws *recws.RecConn
}

const Name = "upbit"

type TickerResponse struct {
	Type string `json:"type"`
	Code string `json:"code"`

	TradePrice     float64 `json:"trade_price"`
	TradeVolume    float64 `json:"trade_volume"`
	TradeTimestamp int64   `json:"trade_timestamp"`

	// Optional fields
	OpeningPrice       float64 `json:"opening_price"`
	HighPrice          float64 `json:"high_price"`
	LowPrice           float64 `json:"low_price"`
	PrevClosingPrice   float64 `json:"prev_closing_price"`
	AccTradePrice      float64 `json:"acc_trade_price"`
	Change             string  `json:"change"`
	ChangePrice        float64 `json:"change_price"`
	SignedChangePrice  float64 `json:"signed_change_price"`
	ChangeRate         float64 `json:"change_rate"`
	SignedChangeRate   float64 `json:"signed_change_rate"`
	AskBid             string  `json:"ask_bid"`
	AccTradeVolume     float64 `json:"acc_trade_volume"`
	TradeDate          string  `json:"trade_date"`
	TradeTime          string  `json:"trade_time"`
	AccAskVolume       float64 `json:"acc_ask_volume"`
	AccBidVolume       float64 `json:"acc_bid_volume"`
	Highest52WeekPrice float64 `json:"highest_52_week_price"`
	Highest52WeekDate  string  `json:"highest_52_week_date"`
	Lowest52WeekPrice  float64 `json:"lowest_52_week_price"`
	Lowest52WeekDate   string  `json:"lowest_52_week_date"`
	MarketState        string  `json:"market_state"`
	IsTradingSuspended bool    `json:"is_trading_suspended"`
	DelistingDate      any     `json:"delisting_date"`
	MarketWarning      string  `json:"market_warning"`
	Timestamp          int64   `json:"timestamp"`
	AccTradePrice24H   float64 `json:"acc_trade_price_24h"`
	AccTradeVolume24H  float64 `json:"acc_trade_volume_24h"`
	StreamType         string  `json:"stream_type"`
}
