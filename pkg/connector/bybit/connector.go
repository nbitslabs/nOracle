package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/recws-org/recws"
)

const Name = "bybit"

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

	wsUrlWithChannels := fmt.Sprintf("%s/v5/public/spot", wsUrl)

	ws := &recws.RecConn{
		KeepAliveTimeout: 10 * time.Second,
	}
	ws.SubscribeHandler = func() error {
		args := make([]string, 0, len(pairs))
		for _, pair := range pairs {
			args = append(args, fmt.Sprintf("tickers.%s", strings.ToUpper(pair)))
		}

		req := SubscriptionMessage{
			Op:    "subscribe",
			ReqId: uuid.New().String(),
			Args:  args,
		}
		return ws.WriteJSON(req)
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

				var tickerResponse TickerResponse
				if err := json.Unmarshal(message, &tickerResponse); err != nil {
					slog.Warn("error unmarshalling message", "error", err, "exchange", Name)
					continue
				}

				out <- connector.TickerUpdate{
					Exchange:  Name,
					Symbol:    tickerResponse.Data.Symbol,
					Price:     tickerResponse.Data.LastPrice,
					Volume:    tickerResponse.Data.Volume24H,
					Timestamp: tickerResponse.Ts,
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
