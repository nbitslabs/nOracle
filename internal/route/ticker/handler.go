package ticker

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/storage"
)

type API struct {
	ctx context.Context

	store storage.Operations[connector.TickerUpdate]
	out   chan connector.TickerUpdate
}

func NewAPI(ctx context.Context, store storage.Operations[connector.TickerUpdate], out chan connector.TickerUpdate) *API {
	api := &API{
		ctx:   ctx,
		store: store,
		out:   out,
	}

	go api.manageStore()

	return api
}

func (a *API) Routes(r *gin.Engine) {
	r.GET("/ticker/:trading/:exchange/:symbol", a.GetTicker)
	r.GET("/price/:trading/:method/:symbol", a.GetPrice)
}

func (a *API) GetTicker(c *gin.Context) {
	exchange := c.Param("exchange")
	symbol := c.Param("symbol")
	trading := c.Param("trading")

	key := strings.ToLower(fmt.Sprintf("%s:%s:%s", exchange, symbol, trading))
	ticker, err := a.store.Get(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticker)
}

func (a *API) GetPrice(c *gin.Context) {
	trading := c.Param("trading")
	method := c.Param("method")
	symbol := c.Param("symbol")
	exchanges := c.QueryArray("exchange")

	if len(exchanges) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No exchanges provided"})
		return
	}

	var price *big.Float
	var err error
	switch method {
	case "average":
		price, err = a.averagePrice(symbol, exchanges, trading)
	case "median":
		price, err = a.medianPrice(symbol, exchanges, trading)
	case "min":
		price, err = a.minPrice(symbol, exchanges, trading)
	case "max":
		price, err = a.maxPrice(symbol, exchanges, trading)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid method"})
		return
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"price": price, "method": method, "symbol": symbol, "exchanges": exchanges})
}

func (a *API) manageStore() {
	for {
		select {
		case ticker := <-a.out:
			switch {
			case ticker.Spot != nil:
				key := strings.ToLower(fmt.Sprintf("%s:%s:spot", ticker.Exchange, ticker.Symbol))
				if err := a.store.Store(key, ticker); err != nil {
					slog.Error("Failed to store ticker", "error", err)
				}
			case ticker.Futures != nil:
				key := strings.ToLower(fmt.Sprintf("%s:%s:futures", ticker.Exchange, ticker.Symbol))
				if err := a.store.Store(key, ticker); err != nil {
					slog.Error("Failed to store ticker", "error", err)
				}
			default:
				slog.Error("Received unknown ticker", "ticker", ticker)
			}
		case <-a.ctx.Done():
			slog.Info("Stopping ticker store manager")
			return
		}
	}
}

func (a *API) averagePrice(symbol string, exchanges []string, trading string) (*big.Float, error) {
	total := big.NewFloat(0)
	count := 0
	for _, exchange := range exchanges {
		ticker, err := a.store.Get(fmt.Sprintf("%s:%s:%s", exchange, symbol, trading))
		if err != nil {
			return nil, fmt.Errorf("ticker not found: %w", err)
		}

		if trading == "spot" {
			total.Add(total, ticker.Spot.Price)
			count++
		} else {
			total.Add(total, ticker.Futures.IndexPrice)
			count++
		}
	}

	average := total.Quo(total, big.NewFloat(float64(count)))

	return average, nil
}

func (a *API) medianPrice(symbol string, exchanges []string, trading string) (*big.Float, error) {
	prices := []*big.Float{}
	for _, exchange := range exchanges {
		ticker, err := a.store.Get(fmt.Sprintf("%s:%s:%s", exchange, symbol, trading))
		if err != nil {
			return nil, fmt.Errorf("ticker not found: %w", err)
		}

		if trading == "spot" {
			prices = append(prices, ticker.Spot.Price)
		} else {
			prices = append(prices, ticker.Futures.IndexPrice)
		}
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Cmp(prices[j]) < 0
	})

	mid := len(prices) / 2
	if len(prices)%2 == 1 {
		return prices[mid], nil
	} else {
		sum := new(big.Float).Add(prices[mid-1], prices[mid])
		return new(big.Float).Quo(sum, big.NewFloat(2)), nil
	}
}

func (a *API) minPrice(symbol string, exchanges []string, trading string) (*big.Float, error) {
	min := big.NewFloat(math.MaxFloat64)

	isSpot := trading == "spot"
	for _, exchange := range exchanges {
		ticker, err := a.store.Get(fmt.Sprintf("%s:%s:%s", exchange, symbol, trading))
		if err != nil {
			return nil, fmt.Errorf("ticker not found: %w", err)
		}

		if isSpot {
			if ticker.Spot.Price.Cmp(min) < 0 {
				min = ticker.Spot.Price
			}
		} else {
			if ticker.Futures.IndexPrice.Cmp(min) < 0 {
				min = ticker.Futures.IndexPrice
			}
		}
	}

	return min, nil
}

func (a *API) maxPrice(symbol string, exchanges []string, trading string) (*big.Float, error) {
	max := big.NewFloat(math.SmallestNonzeroFloat64)

	isSpot := trading == "spot"
	for _, exchange := range exchanges {
		ticker, err := a.store.Get(fmt.Sprintf("%s:%s:%s", exchange, symbol, trading))
		if err != nil {
			return nil, fmt.Errorf("ticker not found: %w", err)
		}

		if isSpot {
			if ticker.Spot.Price.Cmp(max) > 0 {
				max = ticker.Spot.Price
			}
		} else {
			if ticker.Futures.IndexPrice.Cmp(max) > 0 {
				max = ticker.Futures.IndexPrice
			}
		}
	}

	return max, nil
}
