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


//This function
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
	//Last_trade float64
	//High float64
	//Low float64
	//Avg float64
	//Vol float64
	//Vol_curr float64
	//Updated int
}


type pairList map[string]Pair


type chain struct {
	multiply float64
	kProfit  float64
	nextLink int8
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


func getInd (c []string, currencies []string) (int8, int8) {
	var i, j int8
	k := int8(0)
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
	nCurrency := len(*t)
	for k := 0; k < nCurrency; k++ {
		for i := 0; i < nCurrency; i++ {
			for j := 0; j < nCurrency; j++ {
				IJ := &(*t)[i][j]
				if (IJ.kProfit != 0) && (i != j) && (i != k) && (j != k) {
					IK := &(*t)[i][k]
					KJ := &(*t)[k][j]
					if (IK.multiply * KJ.multiply) > IJ.multiply {
						IJ.multiply = IK.multiply * KJ.multiply * (1 - fee)
						IJ.nextLink = IK.nextLink
					}
				}

			}
		}
	}
}

func main() {
	app := new(App)
	currencies, _ := app.getCurrencies()
	pairs, _ := app.getPairs()
	table := pairs.generateTable(currencies)
	table.FloydWarshall()

	nCurrencies := int8(len(currencies))

	for i := int8(0); i < nCurrencies; i++ {
		for j := int8(0); j < nCurrencies; j++ {
			if (i != j) && (table[i][j].multiply * table[i][j].kProfit > 1) {
				fmt.Println(table[i][j].multiply * table[i][j].kProfit)
				fmt.Print(currencies[i], "-->")
				for k := table[i][j].nextLink; k != j; {
					fmt.Print(currencies[k], "-->")
					k = table[k][j].nextLink
				}
				fmt.Print(currencies[j])
			}
		}
	}

}
