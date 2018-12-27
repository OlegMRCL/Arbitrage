package exchange

type Provider interface {
	getCurrencies() ([]string, error)
	getPairs() (PairList, error)
	getPriceTable([]string) PriceTable
	getFee() float64
	NewExchange() Exchange
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
