package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/recws-org/recws"
)

func NewConnector(ctx context.Context, wsUrl string, pairs []string) (connector.ExchangeConnector, error) {
	if wsUrl == "" {
		return nil, fmt.Errorf("wsUrl is required")
	}
	if len(pairs) == 0 {
		return nil, fmt.Errorf("pairs are required")
	}

	channels := []string{}
	for _, pair := range pairs {
		channels = append(channels, fmt.Sprintf("%s@ticker", strings.ToLower(pair)))
	}

	wsUrlWithChannels := fmt.Sprintf("%s/stream?streams=%s", wsUrl, strings.Join(channels, "/"))

	ws := &recws.RecConn{
		KeepAliveTimeout: 10 * time.Second,
	}
	ws.Dial(wsUrlWithChannels, nil)

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

				var t Ticker
				if err := json.Unmarshal(message, &t); err != nil {
					slog.Warn("error unmarshalling message", "error", err, "exchange", Name)
					continue
				}

				out <- connector.TickerUpdate{
					Exchange:  Name,
					Symbol:    t.Stream.Symbol,
					Price:     t.Stream.LastPrice,
					Volume:    t.Stream.Volume,
					Timestamp: t.Stream.EventTime,
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
