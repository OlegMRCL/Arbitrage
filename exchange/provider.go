package exchange

import (
	"ArbitrageFinder/arbitrage"
	"ArbitrageFinder/exchange/exmo"
)

type Provider interface {
	GetCurrencies() ([]string, error)
	GetPairs() (arbitrage.PairList, error)
}


func GetProvider(pType ProviderType) Provider {
	switch pType {
	case ExmoProvider:
		return new(exmo.Exmo)
	case BinanceProvider:
		return nil //new(binanceProvider.Binance)
	}
	return nil
}


type ProviderType string


const (
	ExmoProvider    ProviderType = "EXMO"
	BinanceProvider ProviderType = "BINANCE"
)


