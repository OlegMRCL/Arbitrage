package main

const fee = 0.0005


type Provider interface {
	GetCurrencies() ([]string, error)
	GetPairs() (PairList, error)
}

type providerType string

func getProvider(pType providerType) Provider {
	switch pType {
	case "EXMO":
		return new(Exmo)
	case "BINANCE":
		return nil //new(binance.Binance)
	}
	return nil
}



func main() {
	provider := getProvider("EXMO")
	arbitrage := NewArbitrage(provider)

	//Get data from EXMO

	//Find arbitrage using the data obtained

	//Show all found chains with profit or message about their absence

}
