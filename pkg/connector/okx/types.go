package okx

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

const Name = "okx"

type SubscriptionMessage struct {
	Op   string        `json:"op"`
	Args []ChannelArgs `json:"args"`
}

type ChannelArgs struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

type TickerData struct {
	InstType  string     `json:"instType"`
	InstID    string     `json:"instId"`
	Last      *big.Float `json:"last"`
	LastSz    *big.Float `json:"lastSz"`
	AskPx     *big.Float `json:"askPx"`
	AskSz     *big.Float `json:"askSz"`
	BidPx     *big.Float `json:"bidPx"`
	BidSz     *big.Float `json:"bidSz"`
	Open24H   *big.Float `json:"open24h"`
	High24H   *big.Float `json:"high24h"`
	Low24H    *big.Float `json:"low24h"`
	VolCcy24H *big.Float `json:"volCcy24h"`
	Vol24H    *big.Float `json:"vol24h"`
	SodUtc0   *big.Float `json:"sodUtc0"`
	SodUtc8   *big.Float `json:"sodUtc8"`
	Ts        string     `json:"ts"`
}

type TickerResponse struct {
	Arg  ChannelArgs  `json:"arg"`
	Data []TickerData `json:"data"`
}
