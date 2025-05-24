package ticker

import (
	"fmt"
	"strings"
)

func HyphenSeparatedToStandardTicker(ticker string) string {
	parts := strings.Split(ticker, "-")
	if len(parts) != 2 {
		// If the ticker is not in the format BTC-USDT, we return the same string
		return ticker
	}

	return fmt.Sprintf("%s%s", strings.ToUpper(parts[0]), strings.ToUpper(parts[1]))
}
