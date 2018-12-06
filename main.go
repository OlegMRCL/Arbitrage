package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strconv"
)

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
		Bid string `json:"Buy_price"`;
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


func (pl *pairList) findPair(currency1 string, currency2 string) Pair {
	if _,ok := (*pl)[currency1+"_"+currency2]; ok {
		return (*pl)[currency1+"_"+currency2]
	}else if _,ok := (*pl)[currency2+"_"+currency1]; ok{
		return (*pl)[currency2+"_"+currency1]
	}else{
		return Pair{0, 0}
	}
}


type profit struct {
	multiply float64
	profit float64
	nextLink string
}



func FloydWarshall(currencies []string) {
	for _, k := range currencies {
		for _, i := range currencies {
			for _, j := range currencies {
				fmt.Println(k, i, j)

			}
		}
	}
}


func main() {
	app := new(App)

	currencies, err := app.getCurrencies()

	pairs, err := app.getPairs()
	fmt.Println(currencies, err)
	fmt.Println(pairs.findPair("XMR", "GUSD"))

}
