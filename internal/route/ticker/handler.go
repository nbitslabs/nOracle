package ticker

import (
	"context"
	"fmt"
	"math/big"
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
	r.GET("/ticker/:trading/:exchange/:symbol", a.GetTicker)
	r.GET("/price/:trading/:method/:symbol", a.GetPrice)
}

func (a *API) GetTicker(c *gin.Context) {
	exchange := strings.ToLower(c.Param("exchange"))
	symbol := strings.ToLower(c.Param("symbol"))
	trading := strings.ToLower(c.Param("trading"))

	key := strings.ToLower(fmt.Sprintf("%s:%s:%s", exchange, symbol, trading))
	ticker, err := a.store.Get(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticker)
}

func (a *API) GetPrice(c *gin.Context) {
	trading := strings.ToLower(c.Param("trading"))
	method := strings.ToLower(c.Param("method"))
	symbol := strings.ToLower(c.Param("symbol"))
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
