package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
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
	result := make(map[string]Pair)
	err := json.Unmarshal(body, &result)
	return result, err
}



type Pair struct {
	Buy_price  string
	Sell_price string
	//Last_trade string
	//High string
	//Low string
	//Avg string
	//Vol string
	//Vol_curr string
	//Updated int
}



func main() {

}
