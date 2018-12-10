package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const fee = 0.002


type App struct {
	http.Client
}


func (app *App) getCurrencies() ([]string, error) {
	resp, _ := app.Get("https://api.exmo.com/v1/currency/")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var currencies []string
	err := json.Unmarshal(body, &currencies)
	return currencies, err
}



func (app *App) getPairs() (pairList, error) {

	resp, _ := app.Get("https://api.exmo.com/v1/ticker/")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	result := make(map[string]struct{
		Bid string `json:"Buy_price"`
		Ask string `json:"Sell_price"`
		})
	err := json.Unmarshal(body, &result)

	pairs := make (pairList)
	for key, value := range result {
		p := new(Pair)
		p.Bid, _ = strconv.ParseFloat(value.Bid, 64)
		p.Ask, _ = strconv.ParseFloat(value.Ask, 64)
		pairs[key] = *p
	}

	return pairs, err
}


type Pair struct {
	Bid  float64
	Ask float64
}


type pairList map[string]Pair


type chain struct {
	product  float64
	kProfit  float64
	nextLink uint8
}


func (pl *pairList)  generateTable (currencies []string) Table {
	nCurrency := len(currencies)
	t := make(Table, nCurrency)
	for k := 0; k < nCurrency; k++ {
		t[k] = make([]chain, nCurrency)
	}
	for k, v := range *pl {
		c := strings.Split(k, "_")
		i, j := getInd(c, currencies)
		t[i][j] = chain{v.Bid * (1 - fee), 1/v.Ask * (1 - fee), j}
		t[j][i] = chain{1/v.Ask * (1 - fee), v.Bid * (1 - fee), i}
	}
	return t
}


func getInd (c []string, currencies []string) (uint8, uint8) {
	var i, j uint8
	k := uint8(0)
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


type Table [][]chain


func (t *Table) FloydWarshall() {
	nCurrency := uint8(len(*t))
	for k := uint8(0); k < nCurrency; k++ {
		for i := uint8(0); i < nCurrency; i++ {
			for j := uint8(0); j < nCurrency; j++ {
				(*t)[i][j] = t.compareChains(i, j, k)
			}
		}
	}
}

func (t *Table) compareChains (i, j, k uint8) chain{
	IJ := &(*t)[i][j]

	if (i != j) && (i != k) && (j != k) {
		IK := &(*t)[i][k]
		KJ := &(*t)[k][j]
		if (IK.product * KJ.product) > IJ.product {
			IJ.product = IK.product * KJ.product * (1 - fee)
			IJ.nextLink = IK.nextLink
		}
	}
	return *IJ
}


func showChains (t *Table, currencies []string) {
	nCurrencies := uint8(len(currencies))
	var count uint8 = 0
	for i := uint8(0); i < nCurrencies; i++ {
		for j := uint8(0); j < nCurrencies; j++ {
			if (i != j) && ((*t)[i][j].product * (*t)[i][j].kProfit > 1) {
				fmt.Println((*t)[i][j].product * (*t)[i][j].kProfit)
				fmt.Print(currencies[i], "-->")
				for k := (*t)[i][j].nextLink; k != j; {
					fmt.Print(currencies[k], "-->")
					k = (*t)[k][j].nextLink
				}
				fmt.Print(currencies[j], "-->", currencies[i], "\n")
				count++
			}
		}
	}
	if count == 0 {
		fmt.Println("Arbitrage is not found")
	}
}



func main() {
	app := new(App)

	//Get data from EXMO
	currencies, _ := app.getCurrencies()
	pairs, _ := app.getPairs()

	//Find arbitrage using the data obtained
	table := pairs.generateTable(currencies)
	table.FloydWarshall()

	//Show all found chains with profit or message about their absence
	showChains(&table, currencies)
}
