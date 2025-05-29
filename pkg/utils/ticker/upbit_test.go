package ticker

import "testing"

func TestUpbitToStandardTicker(t *testing.T) {
	tests := []struct {
		upbitTicker string
		want        string
	}{
		{"USDT-BTC", "BTCUSDT"},
		{"USDT-ETH", "ETHUSDT"},
		{"", ""},
		{"BTCUSD", "BTCUSD"},
	}

	for _, test := range tests {
		if got := UpbitToStandardTicker(test.upbitTicker); got != test.want {
			t.Errorf("UpbitToStandardTicker(%s) = %s, want %s", test.upbitTicker, got, test.want)
		}
	}
}
