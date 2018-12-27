package main

import (
	"fmt"
	"github.com/OlegMRCL/ArbitrageFinder/exchange"
	"net/http"
)

func main() {

	http.HandleFunc("/exmo", func(w http.ResponseWriter, r *http.Request) {

		provider := exchange.GetProvider(exchange.Exmo)
		exch := provider.NewExchange()

		results := exch.FindArbitrage()

		fmt.Println(len(results), "chains were found \n")
		fmt.Fprintln(w, len(results), "chains were found \n")

		for k, v := range results {
			fmt.Fprintln(w, "Profit: ", v, "\n", k, "\n")
		}
	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)

}
