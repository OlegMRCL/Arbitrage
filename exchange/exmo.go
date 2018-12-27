package exchange

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ExmoProvider struct {
	http.Client
}

const (
	host             = "https://api.exmo.com"
	endpointCurrency = "/v1/currency/"
	endpointTicker   = "/v1/ticker/"

	Fee = 0.002
)

//Список валютных пар
type PairList map[string]Pair

//Данные о валютной паре
type Pair struct {
	Bid float64 //цена спроса
	Ask float64 //цена предложения
}

//Возвращает объект типа Exchange с полями, заполненными данными
func (e *ExmoProvider) NewExchange() Exchange {

	currencies, _ := e.getCurrencies()

	priceTable := e.getPriceTable(currencies)

	fee := e.getFee()

	return Exchange{
		Currencies: currencies,
		PriceTable: priceTable,
		Fee:        fee,
	}
}

//Возвращает список валют биржи
func (e *ExmoProvider) getCurrencies() (currencies []string, err error) {
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

//Возвращает заполненную матрицу PriceTable
func (e *ExmoProvider) getPriceTable(currencies []string) PriceTable {
	pairs, _ := e.getPairs()

	numCurr := len(currencies)

	pt := make(PriceTable, numCurr)
	for k := 0; k < numCurr; k++ {
		pt[k] = make([]float64, numCurr)
	}

	for k, v := range pairs {
		c := strings.Split(k, "_")
		i, j := getInd(c, currencies)

		pt[i][j] = v.Bid
		pt[j][i] = 1 / v.Ask
	}

	return pt
}

//Возвращает список валютных пар с ценами спроса и предложения
func (e *ExmoProvider) getPairs() (pairs PairList, err error) {

	pairs = make(PairList)
	resp, err := e.Get(host + endpointTicker)

	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			result := make(map[string]struct {
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
func (e *ExmoProvider) getFee() float64 {
	return Fee
}

//Принимает срез с названиями валют,
// возвращает их индексы в срезе e.Currencies
// (только для двух первых валют в передаваемом срезе!)
func getInd(c []string, currencies []string) (i, j int) {
	k := 0
	for foundInd := 0; foundInd < 2; {

		switch {
		case currencies[k] == c[0]:
			i = k
			foundInd++
		case currencies[k] == c[1]:
			j = k
			foundInd++
		}
		k++

	}
	return i, j
}
