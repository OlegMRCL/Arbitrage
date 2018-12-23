package exchange

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type ExmoProvider struct {
	http.Client
}

const (
	host = "https://api.exmo.com"
	endpointCurrency = "/v1/currency/"
	endpointTicker = "/v1/ticker/"

	Fee = 0.002
)


//Возвращает список валют биржи
func (e *ExmoProvider) GetCurrencies() (currencies []string, err error) {
	resp, err := e.Get(host + endpointCurrency)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			err = json.Unmarshal(body, &currencies)
		}
	}
	return
}


//Возвращает список валютных пар с ценами спроса и предложения
func (e *ExmoProvider) GetPairs() (pairs PairList, err error) {

	pairs = make (PairList)
	resp, err := e.Get(host + endpointTicker)

	if err == nil{
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			result := make(map[string]struct{
				Bid string `json:"Buy_price"`
				Ask string `json:"Sell_price"`
			})

			err = json.Unmarshal(body, &result)

			for key, value := range result {
				p := new(Pair)
				p.Bid, _ = strconv.ParseFloat(value.Bid, 64)
				p.Ask, _ = strconv.ParseFloat(value.Ask, 64)
				pairs[key] = *p
			}
		}
	}
	return
}


//Возвращает действующую на бирже комиссию
func (e *ExmoProvider) GetFee() (float64) {
	return Fee
}
