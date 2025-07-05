package ticker

import (
	"fmt"
	"log/slog"
	"strings"
)

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
