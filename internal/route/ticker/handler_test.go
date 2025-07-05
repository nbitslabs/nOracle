package ticker

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	c := make(chan connector.TickerUpdate)
	api := NewAPI(
		context.Background(),
		storage.NewMemory[connector.TickerUpdate](),
		c,
	)

	router := gin.Default()
	api.Routes(router)

	data := connector.TickerUpdate{
		Exchange: "binance",
		Symbol:   "BTCUSDT",
		Spot:     &connector.SpotPriceUpdate{Price: big.NewFloat(100000)},
	}
	c <- data
	// Delay to ensure the ticker is stored
	time.Sleep(200 * time.Millisecond)

	t.Run("GET /ticker/:trading/:exchange/:symbol", func(t *testing.T) {
		t.Run("should return 200", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/ticker/spot/binance/BTCUSDT", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			var response connector.TickerUpdate
			jsonErr := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, jsonErr)

			assert.Equal(t, data.Exchange, response.Exchange)
			assert.Equal(t, data.Symbol, response.Symbol)
			assert.Equal(t, data.Spot.Price.String(), response.Spot.Price.String())
		})

		t.Run("should return 404 when ticker is not found", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/ticker/spot/binance/ETHUSDT", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("should return 404 when trading is not found", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/ticker/futures/binance/BTCUSDT", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("should return 404 when exchange is not found", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/ticker/spot/bybit/BTCUSDT", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})
	})

	t.Run("GET /price/:trading/:method/:symbol", func(t *testing.T) {
		t.Run("should return 200", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/price/spot/average/BTCUSDT?exchange=binance", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusOK, recorder.Code)

			var response map[string]interface{}
			jsonErr := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, jsonErr)

			assert.Equal(t, "100000", response["price"])
		})

		t.Run("should return 400 when no exchanges are provided", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/price/spot/average/BTCUSDT", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("should return 400 when method is invalid", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/price/spot/invalid/BTCUSDT?exchange=binance", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("should return 404 when ticker is not found", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/price/spot/average/ETHUSDT?exchange=binance", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})
	})
}
