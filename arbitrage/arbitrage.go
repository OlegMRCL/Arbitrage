package arbitrage

import (
	"fmt"
	"strings"
	"ArbitrageFinder/exchange"
)


type Arbitrage struct {
	Provider   exchange.Provider
	Currencies []string
	PairList   PairList
	Table      Table
}


func NewArbitrage(p exchange.Provider) (Arbitrage) {
	a := Arbitrage {
		Provider: p,
	}
	a.Currencies, _ = a.Provider.GetCurrencies()
	a.PairList, _ = a.Provider.GetPairs()
	a.Table = a.GenerateTable()
	a.FloydWarshall()

	return a
}


type Pair struct {
	Bid  float64
	Ask float64
}


type PairList map[string]Pair


type Table [][]chain


type chain struct {
	product  float64
	kProfit  float64
	nextLink uint8
	visited map[uint8]bool
}


func (a *Arbitrage) GenerateTable() Table {
	numCurrencies := len(a.Currencies)
	t := make(Table, numCurrencies)
	for k := 0; k < numCurrencies; k++ {
		t[k] = make([]chain, numCurrencies)
	}
	for k, v := range a.PairList {
		c := strings.Split(k, "_")
		i, j := a.getInd(c)

		t[i][j] = chain{
			product: v.Bid * (1 - Fee),
			kProfit: 1/v.Ask * (1 - fee),
			nextLink: j,
			visited: map[uint8]bool {i: true},
		}

		t[j][i] = chain{
			product: 1/v.Ask * (1 - fee),
			kProfit: v.Bid * (1 - fee),
			nextLink: i,
			visited: map[uint8]bool {j: true},
		}

	}
	return t
}


func (a *Arbitrage) getInd (c []string) (uint8, uint8) {
	var i, j uint8
	k := uint8(0)
	for foundInd := 0; foundInd < 2; {
		switch {
		case a.Currencies[k] == c[0]:
			i = k
			foundInd++
		case a.Currencies[k] == c[1]:
			j = k
			foundInd++
		}
		k++
	}
	return i, j
}


func (a *Arbitrage) FloydWarshall() {
	nCurrency := uint8(len(a.Currencies))
	var k, i, j uint8
	for k = 0; k < nCurrency; k++ {
		for i = 0; i < nCurrency; i++ {
			for j = 0; j < nCurrency; j++ {
				a.Table[i][j] = a.compareChains(i, j, k)
			}
		}
	}
}


func (a *Arbitrage) compareChains (i, j, k uint8) chain{
	IJ := &a.Table[i][j]
	if (IJ.kProfit != 0) && (i != j) && (i != k) && (j != k) {
		IK := &a.Table[i][k]
		KJ := &a.Table[k][j]
		if ((IK.product * KJ.product) > IJ.product) && (withoutRepeats(*IK, *KJ)) {
			IJ.product = IK.product * KJ.product * (1 - fee)
			IJ.nextLink = IK.nextLink
			IJ.visited = combineVisited(*IK, *KJ)
		}
	}
	return *IJ
}


func withoutRepeats(IK, KJ chain) bool{
	for k, v := range IK.visited {
		if v && KJ.visited[k] {
			return false
		}
	}
	return true
}


func combineVisited (IK, KJ chain) map[uint8]bool {
	for k := range KJ.visited {
		IK.visited[k] = true
	}
	return IK.visited
}

func (a *Arbitrage) calculateChains() {

}

func (a *Arbitrage) PrintChains() {
	numCurrencies := uint8(len(a.Currencies))
	var count uint8
	var i, j uint8
	for i = 0; i < numCurrencies; i++ {
		for j = 0; j < numCurrencies; j++ {
			if (i != j) && (a.Table[i][j].product * a.Table[i][j].kProfit > 1) {
				fmt.Println(a.Table[i][j].product * a.Table[i][j].kProfit)
				fmt.Print(a.Currencies[i], "-->")
				k := a.Table[i][j].nextLink
				for  k != j {
					fmt.Print(a.Currencies[k], "-->")
					k = a.Table[k][j].nextLink
				}
				fmt.Print(a.Currencies[j], "-->", a.Currencies[i], "\n")
				count++
			}
		}
	}
	if count == 0 {
		fmt.Println("ArbitrageFinder is not found")
	}
}
