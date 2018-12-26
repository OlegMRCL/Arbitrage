package main

import (
	"github.com/OlegMRCL/ArbitrageFinder/exchange"
	"fmt"
	"net/http"
)


func main() {



	http.HandleFunc("/exmo", func(w http.ResponseWriter, r *http.Request){

		provider := exchange.GetProvider(exchange.Exmo)

		exch := exchange.NewExchange(provider)

		pt := exch.GetPriceTable()

		at := exch.FindArbitrage(pt)

		results := exch.GetProfitableChains(at, pt)

		fmt.Println(results)

		if len(results) == 0 {
			fmt.Fprint(w,"Arbitrage is not found!")
		} else {
			for _, v := range results {
				fmt.Fprint(w,"Profit: ", v.Profit, "\n")
				for _, i := range v.Path {
					fmt.Fprint(w, exch.Currencies[i], "-->")
				}
				fmt.Fprint(w, exch.Currencies[v.Path[0]], "\n")
			}
		}

	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)

}

