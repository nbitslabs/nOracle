package ticker

func CoinbaseToStandardTicker(ticker string) string {
	// Coinbase ticker format is like this: BTC-USD
	// We need to convert it to the standard format: BTCUSD
	// We need to remove the dash and convert the uppercase to lowercase

	return HyphenSeparatedToStandardTicker(ticker)
}
