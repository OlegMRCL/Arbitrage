package main

import (
	"github.com/OlegMRCL/ArbitrageFinder/exchange"
	"fmt"
	"net/http"
)


func main() {



	http.HandleFunc("/exmo", func(w http.ResponseWriter, r *http.Request){

		provider := exchange.GetProvider(exchange.Exmo)

		exch := provider.NewExchange()

		results := exch.FindArbitrage()

		fmt.Println(results)

		if len(results) == 0 {
			fmt.Fprint(w,"Arbitrage is not found!")
		} else {
			for k, v := range results {
				fmt.Fprint(w,"Profit: ", v, "\n", k)
			}
		}

	})

	fmt.Println("Server is listening...")
	http.ListenAndServe("localhost:8181", nil)

}

