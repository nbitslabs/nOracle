package okx

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/utils/ticker"
)

func NewConnector(ctx context.Context, wsUrl string, pairs []string) (connector.ExchangeConnector, error) {
	if wsUrl == "" {
		return nil, fmt.Errorf("wsUrl is required")
	}
	if len(pairs) == 0 {
		return nil, fmt.Errorf("pairs are required")
	}

	wsUrlWithChannels := fmt.Sprintf("%s/ws/v5/public", wsUrl)

	ws, _, err := websocket.DefaultDialer.Dial(wsUrlWithChannels, nil)
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		req := SubscriptionMessage{
			Op:   "subscribe",
			Args: []ChannelArgs{{Channel: "tickers", InstId: pair}},
		}

		if err := ws.WriteJSON(req); err != nil {
			return nil, err
		}
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
		return c.ws.Close()
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

				if len(tickerResponse.Data) != 1 {
					continue
				}

				data := tickerResponse.Data[0]

				ts, err := strconv.ParseInt(data.Ts, 10, 64)
				if err != nil {
					slog.Warn("error parsing timestamp", "error", err, "exchange", Name)
					continue
				}

				out <- connector.TickerUpdate{
					Exchange:  Name,
					Symbol:    ticker.OKXToStandardTicker(data.InstID),
					Price:     data.Last,
					Volume:    data.LastSz,
					Timestamp: ts,
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
