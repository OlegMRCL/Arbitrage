package main

import (
	"ArbitrageFinder/arbitrage"
	"ArbitrageFinder/exchange"
)

func main() {
	provider := exchange.GetProvider(exchange.ExmoProvider)
	finder := arbitrage.NewArbitrage(provider)

	finder.PrintChains()

}
