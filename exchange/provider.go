package exchange

type Provider interface {
	GetCurrencies() ([]string, error)
	GetPairs() (PairList, error)
	GetFee() (float64)
}


type ProviderType string


const (
	Exmo    ProviderType = "EXMO"
	Binance ProviderType = "BINANCE"
)

//В соответствии с выбранной биржей возвращает
func GetProvider(pType ProviderType) Provider {
	switch pType {
	case Exmo:
		return new(ExmoProvider)
	case Binance:
		return nil //new(BinanceProvider)
	}
	return nil
}





