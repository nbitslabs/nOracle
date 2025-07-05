package ticker

import (
	"fmt"
	"math"
	"math/big"
	"sort"
)

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
