package ticker

import (
	"fmt"
	"strings"
)

func OKXToStandardTicker(okxTicker string) string {
	// OKX ticker format is like this: BTC-USDT
	// We need to convert it to the standard format: BTCUSDT
	// We need to remove the dash and convert the uppercase to lowercase

	parts := strings.Split(okxTicker, "-")
	if len(parts) != 2 {
		// If the ticker is not in the format BTC-USDT, we return the same string
		return okxTicker
	}

	return fmt.Sprintf("%s%s", strings.ToUpper(parts[0]), strings.ToUpper(parts[1]))
}
