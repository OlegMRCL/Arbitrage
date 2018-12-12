package provider

import (
	"ArbitrageFinder/provider/exmo"
	"ArbitrageFinder/arbitrage"
)

type Provider interface {
	GetCurrencies() ([]string, error)
	GetPairs() (arbitrage.PairList, error)
}


const (
	exmoProvider    ProviderType = "EXMO"
	binanceProvider ProviderType = "BINANCE"
)


func getProvider(pType ProviderType) Provider {
	switch pType {
	case exmoProvider:
		return new(exmo.Exmo)
	case binanceProvider:
		return nil //new(binanceProvider.Binance)
	}
	return nil
}


type ProviderType string
