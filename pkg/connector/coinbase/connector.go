package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/utils/ticker"
	"github.com/recws-org/recws"
)

const Name = "coinbase"

type Connector struct {
	ctx   context.Context
	pairs []string
	ws    *recws.RecConn
}

func NewConnector(ctx context.Context, wsUrl string, pairs []string) (connector.ExchangeConnector, error) {
	if wsUrl == "" {
		return nil, fmt.Errorf("wsUrl is required")
	}
	if len(pairs) == 0 {
		return nil, fmt.Errorf("pairs are required")
	}

	ws := &recws.RecConn{
		KeepAliveTimeout: 10 * time.Second,
	}
	ws.SubscribeHandler = func() error {
		req := SubscriptionMessage{
			Type:     "subscribe",
			Channels: []Channel{{Name: "ticker", ProductIds: pairs}},
		}

		return ws.WriteJSON(req)
	}
	ws.Dial(wsUrl, nil)

	return &Connector{
		ctx:   ctx,
		pairs: pairs,
		ws:    ws,
	}, nil
}

func (c *Connector) Close() error {
	if c.ctx != nil {
		c.ctx.Done()
	}
	if c.ws != nil {
		c.ws.Close()
	}
	return nil
}

func (c *Connector) StreamTickers(ctx context.Context, out chan<- connector.TickerUpdate) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := c.ws.ReadMessage()
				if err != nil {
					slog.Warn("error reading message", "error", err, "exchange", Name)
					continue
				}

				var tickerResponse TickerResponse
				if err := json.Unmarshal(message, &tickerResponse); err != nil {
					slog.Warn("error unmarshalling message", "error", err, "exchange", Name)
					continue
				}

				out <- connector.TickerUpdate{
					Exchange:  Name,
					Symbol:    ticker.CoinbaseToStandardTicker(tickerResponse.ProductID),
					Price:     tickerResponse.Price,
					Volume:    tickerResponse.Volume24H,
					Timestamp: tickerResponse.Time.Unix(),
				}
			}
		}
	}()
	return nil
}

func (c *Connector) Name() string {
	return Name
}

func (c *Connector) Tickers() []string {
	return c.pairs
}
