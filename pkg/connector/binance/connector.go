package binance

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/nbitslabs/nOracle/pkg/connector"
)

func NewConnector(ctx context.Context, wsUrl string, pairs []string) (connector.ExchangeConnector, error) {
	if wsUrl == "" {
		return nil, fmt.Errorf("wsUrl is required")
	}
	if len(pairs) == 0 {
		return nil, fmt.Errorf("pairs are required")
	}

	return &Connector{
		ctx:   ctx,
		pairs: pairs,
	}, nil
}

func (c *Connector) Close() error {
	if c.ctx != nil {
		c.ctx.Done()
	}

	return nil
}

func (c *Connector) StreamTickers(ctx context.Context, out chan<- connector.TickerUpdate) error {
	go c.streamSpot(ctx, out)
	go c.streamFutures(ctx, out)

	return nil
}

func (c *Connector) Name() string {
	return Name
}

func (c *Connector) Tickers() []string {
	return c.pairs
}

func (c *Connector) streamSpot(ctx context.Context, out chan<- connector.TickerUpdate) error {
	sstream := make(chan *binance.WsMarketStatEvent)
	doneC, stopC, err := binance.WsCombinedMarketStatServe(c.pairs, func(event *binance.WsMarketStatEvent) {
		sstream <- event
	}, func(err error) {
		slog.Error("error", "error", err, "exchange", Name)
	})
	if err != nil {
		return err
	}

S:
	for {
		select {
		case <-doneC:
			break S
		case <-ctx.Done():
			stopC <- struct{}{}
			break S
		// Spot Events
		case event := <-sstream:
			price, ok := big.NewFloat(0).SetString(event.LastPrice)
			if !ok {
				slog.Error("error parsing price", "error", err, "exchange", Name)
				continue
			}
			volume, ok := big.NewFloat(0).SetString(event.BaseVolume)
			if !ok {
				slog.Error("error parsing volume", "error", err, "exchange", Name)
				continue
			}

			out <- connector.TickerUpdate{
				Exchange: Name,
				Symbol:   event.Symbol,
				Spot: &connector.SpotPriceUpdate{
					Price:  price,
					Volume: volume,
				},
				Timestamp: event.Time,
			}
		}
	}

	return nil
}

func (c *Connector) streamFutures(ctx context.Context, out chan<- connector.TickerUpdate) error {
	fstream := make(chan *futures.WsMarkPriceEvent)
	doneC, stopC, err := futures.WsCombinedMarkPriceServe(c.pairs, func(event *futures.WsMarkPriceEvent) {
		fstream <- event
	}, func(err error) {
		slog.Error("error", "error", err, "exchange", Name)
	})
	if err != nil {
		return err
	}

F:
	for {
		select {
		case <-doneC:
			break F
		case <-ctx.Done():
			stopC <- struct{}{}
			break F
		// Futures Events
		case event := <-fstream:
			mp, ok := big.NewFloat(0).SetString(event.MarkPrice)
			if !ok {
				slog.Error("error parsing mark price", "error", err, "exchange", Name)
				continue
			}
			ip, ok := big.NewFloat(0).SetString(event.IndexPrice)
			if !ok {
				slog.Error("error parsing index price", "error", err, "exchange", Name)
				continue
			}
			fr, ok := big.NewFloat(0).SetString(event.FundingRate)
			if !ok {
				slog.Error("error parsing funding rate", "error", err, "exchange", Name)
				continue
			}

			out <- connector.TickerUpdate{
				Exchange: Name,
				Symbol:   event.Symbol,
				Futures: &connector.FuturesPriceUpdate{
					MarkPrice:   mp,
					IndexPrice:  ip,
					FundingRate: fr,
				},
				Timestamp: event.Time,
			}
		}
	}

	return nil
}
