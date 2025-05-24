package ticker

func OKXToStandardTicker(okxTicker string) string {
	// OKX ticker format is like this: BTC-USDT
	// We need to convert it to the standard format: BTCUSDT
	// We need to remove the dash and convert the uppercase to lowercase

	return HyphenSeparatedToStandardTicker(okxTicker)
}
