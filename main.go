package main

import (
	"ArbitrageFinder/arbitrage"
	"ArbitrageFinder/provider/exmo"
)







func main() {
	provider := getProvider("EXMO")
	finder := arbitrage.NewArbitrage(provider)

	finder.printChains()

}
