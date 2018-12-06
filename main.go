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



func (app *App) getPairs() (map[string]Pair, error) {

	resp, _ := app.Get("https://api.exmo.com/v1/ticker/")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	result := make(map[string]struct{Bid string `json:"Buy_price"`; Ask string `json:"Sell_price"`})
	err := json.Unmarshal(body, &result)

	pairs := make (map[string]Pair)
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


func main() {
	app := new(App)
	fmt.Println(app.getPairs())
}
