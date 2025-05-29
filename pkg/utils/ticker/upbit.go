package ticker

import (
	"fmt"
	"strings"
)

func UpbitToStandardTicker(upbitTicker string) string {
	// Upbit ticker format is like this: USDT-BTC
	// We need to convert it to the standard format: USDTBTC
	// We need to remove the dash and convert the uppercase to lowercase

	parts := strings.Split(upbitTicker, "-")
	if len(parts) != 2 {
		// If the ticker is not in the format USDT-BTC, we return the same string
		return upbitTicker
	}

	return fmt.Sprintf("%s%s", strings.ToUpper(parts[1]), strings.ToUpper(parts[0]))
}
