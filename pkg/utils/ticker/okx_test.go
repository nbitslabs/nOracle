package ticker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOKXToStandardTicker(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "BTC-USDT",
			expected: "BTCUSDT",
		},
		{
			input:    "ETH-USDT",
			expected: "ETHUSDT",
		},
		// Returns the same string if it's not a valid OKX ticker
		{
			input:    "BTC-USDT-SWAP",
			expected: "BTC-USDT-SWAP",
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "BTCUSD",
			expected: "BTCUSD",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := OKXToStandardTicker(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
