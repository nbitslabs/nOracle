package info

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nbitslabs/nOracle/internal/oracle"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	api := NewAPI(context.Background(), &oracle.Info{
		AvailableExchanges: []string{"binance", "okx"},
		AvailableSymbols:   map[string][]string{"okx": {"BTCUSDT", "ETHUSDT"}, "binance": {"BTCUSDT", "ETHUSDT"}},
	})

	router := gin.Default()
	api.Routes(router)

	t.Run("GET /info/exchanges", func(t *testing.T) {
		t.Run("should return 200", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/info/exchanges", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusOK, recorder.Code)

			var response map[string][]string
			jsonErr := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, jsonErr)
			assert.Equal(t, []string{"binance", "okx"}, response["exchanges"])
		})
	})

	t.Run("GET /info/symbols", func(t *testing.T) {
		t.Run("should return 200", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/info/symbols", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, http.StatusOK, recorder.Code)

			var response map[string][]string
			jsonErr := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, jsonErr)
			assert.Equal(t, map[string][]string{"okx": {"BTCUSDT", "ETHUSDT"}, "binance": {"BTCUSDT", "ETHUSDT"}}, response)
		})
	})
}
