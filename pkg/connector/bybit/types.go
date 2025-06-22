package bybit

import (
	"math/big"
)

type SubscriptionMessage struct {
	ReqId string   `json:"req_id"`
	Op    string   `json:"op"`
	Args  []string `json:"args"`
}

type SpotTickerResponse struct {
	Topic string         `json:"topic"`
	Ts    int64          `json:"ts"`
	Type  string         `json:"type"`
	Cs    int64          `json:"cs"`
	Data  SpotTickerData `json:"data"`
}

type FuturesTickerResponse struct {
	Topic string            `json:"topic"`
	Ts    int64             `json:"ts"`
	Type  string            `json:"type"`
	Cs    int64             `json:"cs"`
	Data  FuturesTickerData `json:"data"`
}

type SpotTickerData struct {
	Symbol       string     `json:"symbol"`
	LastPrice    *big.Float `json:"lastPrice"`
	HighPrice24H *big.Float `json:"highPrice24h"`
	LowPrice24H  *big.Float `json:"lowPrice24h"`
	PrevPrice24H *big.Float `json:"prevPrice24h"`
	Volume24H    *big.Float `json:"volume24h"`
	Turnover24H  *big.Float `json:"turnover24h"`
	Price24HPcnt *big.Float `json:"price24hPcnt"`
	// UsdIndexPrice *big.Float `json:"usdIndexPrice"` // This is not needed and prone to empty string
}

type FuturesTickerData struct {
	Symbol      string     `json:"symbol"`
	LastPrice   *big.Float `json:"lastPrice,omitempty"`
	MarkPrice   *big.Float `json:"markPrice,omitempty"`
	IndexPrice  *big.Float `json:"indexPrice,omitempty"`
	Volume24H   *big.Float `json:"volume24h,omitempty"`
	FundingRate *big.Float `json:"fundingRate,omitempty"`
}
