package info

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nbitslabs/nOracle/internal/oracle"
)

type API struct {
	ctx  context.Context
	info *oracle.Info
}

func NewAPI(ctx context.Context, info *oracle.Info) *API {
	return &API{
		ctx:  ctx,
		info: info,
	}
}

func (a *API) Routes(r *gin.Engine) {
	r.GET("/info/exchanges", a.GetExchanges)
	r.GET("/info/symbols", a.GetSymbols)
}

func (a *API) GetExchanges(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"exchanges": a.info.AvailableExchanges})
}

func (a *API) GetSymbols(c *gin.Context) {
	c.JSON(http.StatusOK, a.info.AvailableSymbols)
}
