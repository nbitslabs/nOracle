package ticker

import (
	"math/big"
	"testing"

	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestAveragePrice(t *testing.T) {
	t.Run("Spot price", func(t *testing.T) {
		t.Run("BTCUSDT average price", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})
			api.store.Store("okx:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10010),
				},
			})

			price, err := api.averagePrice("BTCUSDT", []string{"binance", "okx"}, "spot")
			assert.NoError(t, err)
			assert.Equal(t, big.NewFloat(10005), price)
		})

		t.Run("Exchange not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})

			price, err := api.averagePrice("BTCUSDT", []string{"binance", "okx", "bybit"}, "spot")
			assert.Error(t, err)
			assert.Nil(t, price)
		})

		t.Run("Symbol not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})
			api.store.Store("okx:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})

			price, err := api.averagePrice("ETHUSDT", []string{"binance", "okx"}, "spot")
			assert.Error(t, err)
			assert.Nil(t, price)
		})
	})

	t.Run("Futures price", func(t *testing.T) {
		t.Run("BTCUSDT futures average price", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})
			api.store.Store("okx:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10010),
				},
			})

			price, err := api.averagePrice("BTCUSDT", []string{"binance", "okx"}, "futures")
			assert.NoError(t, err)
			assert.Equal(t, big.NewFloat(10005), price)
		})

		t.Run("Futures exchange not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})

			price, err := api.averagePrice("BTCUSDT", []string{"binance", "okx", "bybit"}, "futures")
			assert.Error(t, err)
			assert.Nil(t, price)
		})

		t.Run("Futures symbol not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})
			api.store.Store("okx:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})

			price, err := api.averagePrice("ETHUSDT", []string{"binance", "okx"}, "futures")
			assert.Error(t, err)
			assert.Nil(t, price)
		})
	})
}

func TestMedianPrice(t *testing.T) {
	t.Run("Spot price", func(t *testing.T) {
		t.Run("BTCUSDT median price", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})
			api.store.Store("okx:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10010),
				},
			})
			api.store.Store("bybit:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10020),
				},
			})

			price, err := api.medianPrice("BTCUSDT", []string{"binance", "okx", "bybit"}, "spot")
			assert.NoError(t, err)
			assert.Equal(t, big.NewFloat(10010), price)
		})

		t.Run("BTCUSDT median price with even number of exchanges", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})
			api.store.Store("okx:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10010),
				},
			})

			price, err := api.medianPrice("BTCUSDT", []string{"binance", "okx"}, "spot")
			assert.NoError(t, err)
			assert.Equal(t, big.NewFloat(10005), price)
		})

		t.Run("Exchange not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})

			price, err := api.medianPrice("BTCUSDT", []string{"binance", "okx"}, "spot")
			assert.Error(t, err)
			assert.Nil(t, price)
		})

		t.Run("Symbol not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})

			price, err := api.medianPrice("ETHUSDT", []string{"binance"}, "spot")
			assert.Error(t, err)
			assert.Nil(t, price)
		})
	})

	t.Run("Futures price", func(t *testing.T) {
		t.Run("BTCUSDT futures median price", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})
			api.store.Store("okx:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10010),
				},
			})
			api.store.Store("bybit:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10020),
				},
			})

			price, err := api.medianPrice("BTCUSDT", []string{"binance", "okx", "bybit"}, "futures")
			assert.NoError(t, err)
			assert.Equal(t, big.NewFloat(10010), price)
		})

		t.Run("BTCUSDT futures median price with even number of exchanges", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})
			api.store.Store("okx:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10010),
				},
			})

			price, err := api.medianPrice("BTCUSDT", []string{"binance", "okx"}, "futures")
			assert.NoError(t, err)
			assert.Equal(t, big.NewFloat(10005), price)
		})

		t.Run("Futures exchange not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})

			price, err := api.medianPrice("BTCUSDT", []string{"binance", "okx"}, "futures")
			assert.Error(t, err)
			assert.Nil(t, price)
		})

		t.Run("Futures symbol not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})

			price, err := api.medianPrice("ETHUSDT", []string{"binance"}, "futures")
			assert.Error(t, err)
			assert.Nil(t, price)
		})
	})
}

func TestMinPrice(t *testing.T) {
	t.Run("Spot price", func(t *testing.T) {
		t.Run("BTCUSDT min price", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})
			api.store.Store("okx:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10010),
				},
			})

			price, err := api.minPrice("BTCUSDT", []string{"binance", "okx"}, "spot")
			assert.NoError(t, err)
			assert.Equal(t, big.NewFloat(10000), price)
		})

		t.Run("Exchange not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})

			price, err := api.minPrice("BTCUSDT", []string{"binance", "okx"}, "spot")
			assert.Error(t, err)
			assert.Nil(t, price)
		})

		t.Run("Symbol not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})

			price, err := api.minPrice("ETHUSDT", []string{"binance"}, "spot")
			assert.Error(t, err)
			assert.Nil(t, price)
		})
	})

	t.Run("Futures price", func(t *testing.T) {
		t.Run("BTCUSDT futures min price", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})
			api.store.Store("okx:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10010),
				},
			})

			price, err := api.minPrice("BTCUSDT", []string{"binance", "okx"}, "futures")
			assert.NoError(t, err)
			assert.Equal(t, big.NewFloat(10000), price)
		})

		t.Run("Futures exchange not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})

			price, err := api.minPrice("BTCUSDT", []string{"binance", "okx"}, "futures")
			assert.Error(t, err)
			assert.Nil(t, price)
		})

		t.Run("Futures symbol not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})

			price, err := api.minPrice("ETHUSDT", []string{"binance"}, "futures")
			assert.Error(t, err)
			assert.Nil(t, price)
		})
	})
}

func TestMaxPrice(t *testing.T) {
	t.Run("Spot price", func(t *testing.T) {
		t.Run("BTCUSDT max price", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})

			api.store.Store("okx:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10010),
				},
			})

			price, err := api.maxPrice("BTCUSDT", []string{"binance", "okx"}, "spot")
			assert.NoError(t, err)
			assert.Equal(t, big.NewFloat(10010), price)
		})

		t.Run("Exchange not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})

			price, err := api.maxPrice("BTCUSDT", []string{"binance", "okx"}, "spot")
			assert.Error(t, err)
			assert.Nil(t, price)
		})

		t.Run("Symbol not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:spot", connector.TickerUpdate{
				Spot: &connector.SpotPriceUpdate{
					Price: big.NewFloat(10000),
				},
			})

			price, err := api.maxPrice("ETHUSDT", []string{"binance"}, "spot")
			assert.Error(t, err)
			assert.Nil(t, price)
		})
	})

	t.Run("Futures price", func(t *testing.T) {
		t.Run("BTCUSDT futures max price", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})

			api.store.Store("okx:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10010),
				},
			})

			price, err := api.maxPrice("BTCUSDT", []string{"binance", "okx"}, "futures")
			assert.NoError(t, err)
			assert.Equal(t, big.NewFloat(10010), price)
		})

		t.Run("Futures exchange not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})

			price, err := api.maxPrice("BTCUSDT", []string{"binance", "okx"}, "futures")
			assert.Error(t, err)
			assert.Nil(t, price)
		})

		t.Run("Futures symbol not found", func(t *testing.T) {
			memory := storage.NewMemory[connector.TickerUpdate]()
			api := &API{
				store: memory,
			}

			api.store.Store("binance:BTCUSDT:futures", connector.TickerUpdate{
				Futures: &connector.FuturesPriceUpdate{
					IndexPrice: big.NewFloat(10000),
				},
			})

			price, err := api.maxPrice("ETHUSDT", []string{"binance"}, "futures")
			assert.Error(t, err)
			assert.Nil(t, price)
		})
	})
}
