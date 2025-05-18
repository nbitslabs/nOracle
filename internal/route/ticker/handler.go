package ticker

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
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
	r.GET("/ticker/:exchange/:symbol", a.GetTicker)
}

func (a *API) GetTicker(c *gin.Context) {
	exchange := c.Param("exchange")
	symbol := c.Param("symbol")

	key := strings.ToLower(fmt.Sprintf("%s:%s", exchange, symbol))
	ticker, err := a.store.Get(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticker)
}

func (a *API) manageStore() {
	for {
		select {
		case ticker := <-a.out:
			key := strings.ToLower(fmt.Sprintf("%s:%s", ticker.Exchange, ticker.Symbol))
			if err := a.store.Store(key, ticker); err != nil {
				slog.Error("Failed to store ticker", "error", err)
			}
		case <-a.ctx.Done():
			slog.Info("Stopping ticker store manager")
			return
		}
	}
}
