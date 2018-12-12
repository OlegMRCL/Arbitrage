package exmo

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"ArbitrageFinder/arbitrage"
)

type Exmo struct {
	http.Client
}

const (
	host = "https://api.exmo.com"
	endpointCurrency = "/v1/currency/"
	endpointTicker = "/v1/ticker/"

	Fee = 0.002
)



func (e *Exmo) GetCurrencies() ([]string, error) {
	resp, _ := e.Get(host + endpointCurrency)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var currencies []string
	err := json.Unmarshal(body, &currencies)
	return currencies, err
}



func (e *Exmo) GetPairs() (arbitrage.PairList, error) {

	resp, _ := e.Get(host + endpointTicker)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	result := make(map[string]struct{
		Bid string `json:"Buy_price"`
		Ask string `json:"Sell_price"`
	})
	err := json.Unmarshal(body, &result)

	pairs := make (arbitrage.PairList)
	for key, value := range result {
		p := new(arbitrage.Pair)
		p.Bid, _ = strconv.ParseFloat(value.Bid, 64)
		p.Ask, _ = strconv.ParseFloat(value.Ask, 64)
		pairs[key] = *p
	}

	return pairs, err
}
