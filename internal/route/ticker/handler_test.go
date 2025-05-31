package ticker

import (
	"math/big"
	"testing"

	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestAveragePrice(t *testing.T) {
	t.Run("BTCUSDT average price", func(t *testing.T) {
		memory := storage.NewMemory[connector.TickerUpdate]()
		api := &API{
			store: memory,
		}

		api.store.Store("binance:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10000),
		})
		api.store.Store("okx:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10010),
		})

		price, err := api.averagePrice("BTCUSDT", []string{"binance", "okx"})
		assert.NoError(t, err)
		assert.Equal(t, big.NewFloat(10005), price)
	})

	t.Run("Exchange not found", func(t *testing.T) {
		memory := storage.NewMemory[connector.TickerUpdate]()
		api := &API{
			store: memory,
		}

		api.store.Store("binance:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10000),
		})

		price, err := api.averagePrice("BTCUSDT", []string{"binance", "okx"})
		assert.Error(t, err)
		assert.Nil(t, price)
	})

	t.Run("Symbol not found", func(t *testing.T) {
		memory := storage.NewMemory[connector.TickerUpdate]()
		api := &API{
			store: memory,
		}

		api.store.Store("binance:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10000),
		})
		api.store.Store("okx:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10000),
		})

		price, err := api.averagePrice("ETHUSDT", []string{"binance", "okx"})
		assert.Error(t, err)
		assert.Nil(t, price)
	})
}

func TestMedianPrice(t *testing.T) {
	t.Run("BTCUSDT median price", func(t *testing.T) {
		memory := storage.NewMemory[connector.TickerUpdate]()
		api := &API{
			store: memory,
		}

		api.store.Store("binance:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10000),
		})
		api.store.Store("okx:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10010),
		})
		api.store.Store("bybit:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10020),
		})

		price, err := api.medianPrice("BTCUSDT", []string{"binance", "okx", "bybit"})
		assert.NoError(t, err)
		assert.Equal(t, big.NewFloat(10010), price)
	})

	t.Run("BTCUSDT median price with even number of exchanges", func(t *testing.T) {
		memory := storage.NewMemory[connector.TickerUpdate]()
		api := &API{
			store: memory,
		}

		api.store.Store("binance:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10000),
		})
		api.store.Store("okx:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10010),
		})

		price, err := api.medianPrice("BTCUSDT", []string{"binance", "okx"})
		assert.NoError(t, err)
		assert.Equal(t, big.NewFloat(10005), price)
	})

	t.Run("Exchange not found", func(t *testing.T) {
		memory := storage.NewMemory[connector.TickerUpdate]()
		api := &API{
			store: memory,
		}

		api.store.Store("binance:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10000),
		})

		price, err := api.medianPrice("BTCUSDT", []string{"binance", "okx"})
		assert.Error(t, err)
		assert.Nil(t, price)
	})

	t.Run("Symbol not found", func(t *testing.T) {
		memory := storage.NewMemory[connector.TickerUpdate]()
		api := &API{
			store: memory,
		}

		api.store.Store("binance:BTCUSDT", connector.TickerUpdate{
			Price: big.NewFloat(10000),
		})

		price, err := api.medianPrice("ETHUSDT", []string{"binance"})
		assert.Error(t, err)
		assert.Nil(t, price)
	})
}
