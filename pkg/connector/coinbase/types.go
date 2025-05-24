package coinbase

import (
	"math/big"
	"time"
)

type SubscriptionMessage struct {
	Type     string    `json:"type"`
	Channels []Channel `json:"channels"`
}

type Channel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

type TickerResponse struct {
	Type        string     `json:"type"`
	Sequence    int64      `json:"sequence"`
	ProductID   string     `json:"product_id"`
	Price       *big.Float `json:"price"`
	Open24H     *big.Float `json:"open_24h"`
	Volume24H   *big.Float `json:"volume_24h"`
	Low24H      *big.Float `json:"low_24h"`
	High24H     *big.Float `json:"high_24h"`
	Volume30D   *big.Float `json:"volume_30d"`
	BestBid     *big.Float `json:"best_bid"`
	BestBidSize *big.Float `json:"best_bid_size"`
	BestAsk     *big.Float `json:"best_ask"`
	BestAskSize *big.Float `json:"best_ask_size"`
	Side        string     `json:"side"`
	Time        time.Time  `json:"time"`
	TradeID     *big.Int   `json:"trade_id"`
	LastSize    *big.Float `json:"last_size"`
}
