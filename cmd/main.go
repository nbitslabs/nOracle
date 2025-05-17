package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/nbitslabs/nOracle/pkg/connector"
	"github.com/nbitslabs/nOracle/pkg/connector/binance"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	exit := make(chan struct{})

	c, err := binance.NewConnector(context.Background(), "wss://stream.binance.com:443", []connector.Symbol{"BTCUSDT", "ETHUSDT"})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	out := make(chan connector.TickerUpdate)
	go c.StreamTickers(context.Background(), out)

	<-sig
	close(exit)
}
