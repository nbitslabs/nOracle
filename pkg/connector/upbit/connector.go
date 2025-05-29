package upbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/utils/ticker"
	"github.com/recws-org/recws"
)

func NewConnector(ctx context.Context, wsUrl string, pairs []string) (connector.ExchangeConnector, error) {
	if wsUrl == "" {
		return nil, fmt.Errorf("wsUrl is required")
	}
	if len(pairs) == 0 {
		return nil, fmt.Errorf("pairs are required")
	}

	wsUrlWithChannels := fmt.Sprintf("%s/ws/v5/public", wsUrl)
	ws := &recws.RecConn{
		KeepAliveTimeout: 10 * time.Second,
	}
	ws.Dial(wsUrlWithChannels, nil)

	req := []map[string]interface{}{
		{"ticket": uuid.NewString()},
		{"type": "ticker", "codes": pairs, "isOnlyRealtime": true},
		{"format": "DEFAULT"},
	}

	if err := ws.WriteJSON(req); err != nil {
		return nil, err
	}

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
					fmt.Println(string(message))
					slog.Warn("error unmarshalling message", "error", err, "exchange", Name)
					continue
				}

				out <- connector.TickerUpdate{
					Exchange:  Name,
					Symbol:    ticker.UpbitToStandardTicker(tickerResponse.Code),
					Price:     big.NewFloat(tickerResponse.TradePrice),
					Volume:    big.NewFloat(tickerResponse.TradeVolume),
					Timestamp: tickerResponse.TradeTimestamp,
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
