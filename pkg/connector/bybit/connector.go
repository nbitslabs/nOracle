package bybit

import (
	"context"
	"encoding/json"
	"fmt"

	bybit "github.com/bybit-exchange/bybit.go.api"
	"github.com/nbitslabs/nOracle/pkg/connector"
)

const Name = "bybit"

type Connector struct {
	ctx   context.Context
	pairs []string

	args []string

	fCache map[string]*FuturesTickerResponse
}

func NewConnector(ctx context.Context, pairs []string) (connector.ExchangeConnector, error) {
	if len(pairs) == 0 {
		return nil, fmt.Errorf("pairs are required")
	}

	args := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		args = append(args, "tickers."+pair)
	}

	return &Connector{
		ctx:    ctx,
		pairs:  pairs,
		args:   args,
		fCache: make(map[string]*FuturesTickerResponse),
	}, nil
}

func (c *Connector) Close() error {
	if c.ctx != nil {
		c.ctx.Done()
	}

	return nil
}

func (c *Connector) StreamTickers(ctx context.Context, out chan<- connector.TickerUpdate) error {
	if err := c.streamSpot(ctx, out); err != nil {
		return err
	}

	if err := c.streamFutures(ctx, out); err != nil {
		return err
	}

	return nil
}

func (c *Connector) Name() string {
	return Name
}

func (c *Connector) Tickers() []string {
	return c.pairs
}

func (c *Connector) streamSpot(ctx context.Context, out chan<- connector.TickerUpdate) error {
	client := bybit.NewBybitPublicWebSocket("wss://stream.bybit.com/v5/public/spot",
		bybit.MessageHandler(func(message string) error {
			var res SpotTickerResponse
			if err := json.Unmarshal([]byte(message), &res); err != nil {
				return err
			}

			out <- connector.TickerUpdate{
				Exchange: Name,
				Symbol:   res.Data.Symbol,
				Spot: &connector.SpotPriceUpdate{
					Price:  res.Data.LastPrice,
					Volume: res.Data.Volume24H,
				},
				Timestamp: res.Ts,
			}

			return nil
		}))

	ws := client.Connect()
	if _, err := ws.SendSubscription(c.args); err != nil {
		return err
	}

	return nil
}

func (c *Connector) streamFutures(ctx context.Context, out chan<- connector.TickerUpdate) error {
	client := bybit.NewBybitPublicWebSocket("wss://stream.bybit.com/v5/public/linear",
		bybit.MessageHandler(func(message string) error {
			var res FuturesTickerResponse
			if err := json.Unmarshal([]byte(message), &res); err != nil {
				return err
			}

			symbol := res.Data.Symbol
			if _, ok := c.fCache[symbol]; !ok {
				c.fCache[symbol] = &res
			}

			if res.Data.LastPrice != nil {
				c.fCache[symbol].Data.LastPrice = res.Data.LastPrice
			}
			if res.Data.MarkPrice != nil {
				c.fCache[symbol].Data.MarkPrice = res.Data.MarkPrice
			}
			if res.Data.IndexPrice != nil {
				c.fCache[symbol].Data.IndexPrice = res.Data.IndexPrice
			}
			if res.Data.Volume24H != nil {
				c.fCache[symbol].Data.Volume24H = res.Data.Volume24H
			}
			if res.Data.FundingRate != nil {
				c.fCache[symbol].Data.FundingRate = res.Data.FundingRate
			}

			// update cache
			c.fCache[symbol].Cs = res.Cs
			c.fCache[symbol].Ts = res.Ts

			out <- connector.TickerUpdate{
				Exchange: Name,
				Symbol:   symbol,
				Futures: &connector.FuturesPriceUpdate{
					MarkPrice:   c.fCache[symbol].Data.MarkPrice,
					IndexPrice:  c.fCache[symbol].Data.IndexPrice,
					FundingRate: c.fCache[symbol].Data.FundingRate,
				},
				Timestamp: res.Ts,
			}

			return nil
		}))

	ws := client.Connect()
	if _, err := ws.SendSubscription(c.args); err != nil {
		return err
	}

	return nil
}
