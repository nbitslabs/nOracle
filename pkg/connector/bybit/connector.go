package bybit

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hirokisan/bybit/v2"
	"github.com/nbitslabs/nOracle/pkg/connector"
)

const Name = "bybit"

type Connector struct {
	ctx   context.Context
	pairs []string

	tickers []bybit.V5WebsocketPublicTickerParamKey
}

func NewConnector(ctx context.Context, wsUrl string, pairs []string) (connector.ExchangeConnector, error) {
	if wsUrl == "" {
		return nil, fmt.Errorf("wsUrl is required")
	}
	if len(pairs) == 0 {
		return nil, fmt.Errorf("pairs are required")
	}

	tickers := make([]bybit.V5WebsocketPublicTickerParamKey, 0, len(pairs))
	for _, pair := range pairs {
		tickers = append(tickers, bybit.V5WebsocketPublicTickerParamKey{
			Symbol: bybit.SymbolV5(pair),
		})
	}

	return &Connector{
		ctx:     ctx,
		pairs:   pairs,
		tickers: tickers,
	}, nil
}

func (c *Connector) Close() error {
	if c.ctx != nil {
		c.ctx.Done()
	}

	return nil
}

func (c *Connector) StreamTickers(ctx context.Context, out chan<- connector.TickerUpdate) error {
	ws := bybit.NewWebsocketClient()

	spot, err := c.streamSpot(ctx, ws, out)
	if err != nil {
		return err
	}

	futures, err := c.streamFutures(ctx, ws, out)
	if err != nil {
		return err
	}

	executors := []bybit.WebsocketExecutor{
		spot,
		futures,
	}

	ws.Start(ctx, executors)
	return nil
}

func (c *Connector) Name() string {
	return Name
}

func (c *Connector) Tickers() []string {
	return c.pairs
}

func (c *Connector) streamSpot(ctx context.Context, ws *bybit.WebSocketClient, out chan<- connector.TickerUpdate) (bybit.V5WebsocketPublicServiceI, error) {
	svc, err := ws.V5().Public(bybit.CategoryV5Spot)
	if err != nil {
		return nil, err
	}

	if _, err := svc.SubscribeTickers(c.tickers, func(res bybit.V5WebsocketPublicTickerResponse) error {
		data := res.Data.Spot
		if data.Symbol == "" {
			return fmt.Errorf("received empty ticker data")
		}

		price, ok := big.NewFloat(0).SetString(data.LastPrice)
		if !ok {
			return fmt.Errorf("error parsing price")
		}
		volume, ok := big.NewFloat(0).SetString(data.Volume24H)
		if !ok {
			return fmt.Errorf("error parsing volume")
		}

		out <- connector.TickerUpdate{
			Exchange: Name,
			Symbol:   string(data.Symbol),
			Spot: &connector.SpotPriceUpdate{
				Price:  price,
				Volume: volume,
			},
			Timestamp: res.TimeStamp,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return svc, nil
}

func (c *Connector) streamFutures(ctx context.Context, ws *bybit.WebSocketClient, out chan<- connector.TickerUpdate) (bybit.V5WebsocketPublicServiceI, error) {
	svc, err := ws.V5().Public(bybit.CategoryV5Linear)
	if err != nil {
		return nil, err
	}

	if _, err := svc.SubscribeTickers(c.tickers, func(res bybit.V5WebsocketPublicTickerResponse) error {
		data := res.Data.LinearInverse
		if data.Symbol == "" {
			return fmt.Errorf("received empty ticker data")
		}

		mp, ok := big.NewFloat(0).SetString(data.MarkPrice)
		if !ok {
			return fmt.Errorf("error parsing price")
		}

		ip, ok := big.NewFloat(0).SetString(data.IndexPrice)
		if !ok {
			return fmt.Errorf("error parsing index price")
		}

		fr, ok := big.NewFloat(0).SetString(data.FundingRate)
		if !ok {
			return fmt.Errorf("error parsing funding rate")
		}

		out <- connector.TickerUpdate{
			Exchange:  Name,
			Symbol:    string(data.Symbol),
			Timestamp: res.TimeStamp,
			Futures: &connector.FuturesPriceUpdate{
				MarkPrice:   mp,
				IndexPrice:  ip,
				FundingRate: fr,
			},
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return svc, nil
}
