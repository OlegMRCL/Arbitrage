package exchange

import (
	"strings"
)


//Данные о бирже
type Exchange struct {
	Currencies []string		//список валют на бирже
	PairList PairList		//список валютных пар на бирже
	Fee float64				//комиссия, взимаемая на бирже
}


//Список валютных пар
type PairList map[string]Pair


//Данные о валютной паре
type Pair struct {
	Bid float64		//цена спроса
	Ask float64		//цена предложения
}


//Возвращает объект типа Exchange с полями, заполненными данными
func NewExchange(p Provider) Exchange {
	exch := Exchange{}
	currencies, err := p.GetCurrencies()

	if err == nil {
		pairList, err := p.GetPairs()

		if err == nil {
			fee := p.GetFee()

			exch.Currencies = currencies
			exch.PairList = pairList
			exch.Fee = fee
		}
	}
	return exch
}


//В таблице PriceTable храянятся цены валют:
// в ячейке [i][j] указана цена валюты i в валюте j,
// а в ячейке [j][i] - валюты j в валюте i.
type PriceTable [][] float64


//Возвращает заполненную матрицу PriceTable
func (e *Exchange) GetPriceTable() PriceTable {
	numCurr := len(e.Currencies)

	pt := make(PriceTable, numCurr)
	for k := 0; k < numCurr; k++ {
		pt[k] = make([]float64, numCurr)
	}

	for k, v := range e.PairList {
		currs := strings.Split(k, "_")
		i, j := e.getInd(currs)

		pt[i][j] = v.Bid
		pt[j][i] = 1/v.Ask
	}

	return pt
}


//Принимает срез с названиями валют,
// возвращает их индексы в срезе e.Currencies
// (только для двух первых валют в передаваемом срезе!)
func (e *Exchange) getInd (c []string) (i,j int) {
	k := 0
	for foundInd := 0; foundInd < 2; {

		switch {
		case e.Currencies[k] == c[0]:
			i = k
			foundInd++
		case e.Currencies[k] == c[1]:
			j = k
			foundInd++
		}
		k++

	}
	return i, j
}


//В матрице ArbitrageTable каждой ячейке [i][j] соответствует цепочка межвалютных обменов,
// начальной и конечной точкой которой являются валюты i и j
type ArbitrageTable [][]chain


//Данные о цепочке межвалютных обменов
type chain struct {
	product  float64		//произведение всех межвалютных цен в цепочке
	nextLink int			//индекс второй валюты в цепочке
	visited map[int]bool	//список "посещенных" валют цепочки (включая первую, и не включая последнюю)
}


//Создает матрицу ArbitrageTable и заполняет ее начальными chain
func (e *Exchange) generateTable(pt PriceTable) ArbitrageTable{
	numCurr := len(e.Currencies)

	at := make(ArbitrageTable, numCurr)
	var i, j int
	for i = 0; i < numCurr; i++ {

		at[i] = make([]chain, numCurr)
		for j = 0; j < numCurr; j++ {

			at[i][j] = chain{
				product: pt[i][j] * (1 - e.Fee),
				nextLink: j,
				visited: map[int]bool {i: true},
			}
		}
	}

	return at
}


//Создает и возвращает матрицу ArbitrageTable, заполненную цепочками chain с наибольшим значением product.
// Поиск наивыгоднейших цепочек работает на основе видоизмененного
// алгоритма Флойда-Уоршелла (аддитивная группа заменена на мултипликативную)
func (e *Exchange) FindArbitrage(pt PriceTable) (at ArbitrageTable){

	at = e.generateTable(pt)

	numCurr := len(pt)

	for k := 0; k < numCurr; k++ {
		for i := 0; i < numCurr; i++ {
			for j := 0; j < numCurr; j++ {
				at[i][j] = at.bestChain(i, j, k)
			}
		}
	}

	return at
}


//Возвращает цепочку, выбранную как наивыгоднейшую:
// выбирает между исходной цепочкой IJ и предлагаемой IKJ, сочлененной из IK и KJ.
func (at *ArbitrageTable) bestChain(i, j, k int) chain  {
	IJ := (*at)[i][j]

	//if (i != j) && (i != k) && (j != k) {
	IK := (*at)[i][k]
	KJ := (*at)[k][j]

	if noMatchingPoints(IK, KJ) {

		if IJ.product < (IK.product * KJ.product) {
			IJ.product = IK.product * KJ.product
			IJ.nextLink = IK.nextLink
			IJ.visited = joinVisited(IK, KJ)
		}
	}
	//}

	return IJ
}


//Данные о цепочке, признанной выгодной (Profit > 1)
type Result struct {
	Profit float64		//относительная величина профита
	Path []int			//валютный "маршрут"
}


//Ищет в матрице ArbitrageTable цепочки с положительной прибылью
func (e *Exchange) GetProfitableChains(at ArbitrageTable, pt PriceTable) (Results []Result) {

	numCurr := len(e.Currencies)
	for i := 0; i < numCurr; i++ {
		for j := 0; j < numCurr; j++ {

			//Профит (как относитлеьная величина) равен произведению всех цен в цепочке at[i][j],
			// умноженному на "обратную" цену и минус комиссия.
			// Если профит > 1, то заносим в срез результатов поиска данные по профитной цепочке.
			if profit := at[i][j].product * pt[j][i] * (1 - e.Fee); profit > 1 {

				r := Result{
					Profit: profit,
					Path:   at.getPath(i,j),
				}

				Results = append(Results, r)
			}
		}
	}

	return Results
}


//Возвращает срез, состоящий из валют цепочки at[i][j]
func (at *ArbitrageTable) getPath(i,j int) (path []int) {

	chain := (*at)[i][j]
	path = append(path, i)

	for i != j {
		i = chain.nextLink
		path = append(path, i)
		chain = (*at)[i][j]
	}

	return path
}


//Проверяет две цепочки на отсутствие повторяющихся точек-валют в их "маршрутах"
func noMatchingPoints(IK, KJ chain) bool{

	for k, v := range IK.visited {
		if v && KJ.visited[k] {
			return false
		}
	}
	return true
}


//Возвращает объединенный список "посещенных" валют двух цепочек
func joinVisited(IK, KJ chain) map[int]bool {
	for k := range KJ.visited {
		IK.visited[k] = true
	}
	return IK.visited
}







