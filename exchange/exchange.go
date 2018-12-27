package exchange

import "fmt"

//Данные о бирже
type Exchange struct {
	Currencies []string   //список валют на бирже
	PriceTable PriceTable //таблица межвалютных цен
	Fee        float64    //комиссия, взимаемая на бирже
}

//В таблице PriceTable храянятся цены валют:
// в ячейке [i][j] указана цена валюты i в валюте j,
// а в ячейке [j][i] - валюты j в валюте i.
type PriceTable [][]float64

//В матрице ArbitrageTable каждой ячейке [i][j] соответствует цепочка межвалютных обменов,
// начальной и конечной точкой которой являются валюты i и j
type ArbitrageTable [][]chain

//Данные о цепочке межвалютных обменов
type chain struct {
	product float64 //произведение всех межвалютных цен в цепочке
	path    []int   //список индексов валют в цепочке
}

//Создает матрицу ArbitrageTable и заполняет ее начальными цепочками chain
func (e *Exchange) generateTable() ArbitrageTable {
	numCurr := len(e.Currencies)

	at := make(ArbitrageTable, numCurr)
	for i := 0; i < numCurr; i++ {

		at[i] = make([]chain, numCurr)
		for j := 0; j < numCurr; j++ {

			at[i][j] = chain{
				product: e.PriceTable[i][j] * (1 - e.Fee),
				path:    []int{i},
			}
		}
	}

	return at
}

//Создает матрицу ArbitrageTable, заполненную цепочками chain с наибольшим значением product.
// Возвращает список профитных цепочек.
// Поиск наивыгоднейших цепочек работает на основе видоизмененного
// алгоритма Флойда-Уоршелла (аддитивная группа заменена на мултипликативную)
func (e *Exchange) FindArbitrage() (results Results) {

	at := e.generateTable()
	results = make(Results)

	numCurr := len(e.PriceTable)

	for k := 0; k < numCurr; k++ {
		for i := 0; i < numCurr; i++ {
			for j := 0; j < numCurr; j++ {

				at[i][j] = at.checkChains(i, j, k)

				if (i == j) && (at[i][j].product > 1) {

					pathString := ""
					for _, v := range at[i][j].path {
						pathString += e.Currencies[v] + "-->"
					}
					pathString += e.Currencies[i]

					results[pathString] = at[i][j].product
				}
			}
		}
	}

	for i := 0; i < numCurr; i++ {
		fmt.Println(i, e.Currencies[i], at[i][i])
	}

	return results
}

type Results map[string]float64

//Возвращает цепочку, выбранную как наивыгоднейшую:
// выбирает между исходной цепочкой IJ и предлагаемой IKJ, сочлененной из IK и KJ.
func (at *ArbitrageTable) checkChains(i, j, k int) chain {
	IJ := (*at)[i][j]
	IK := (*at)[i][k]
	KJ := (*at)[k][j]

	if noMatchingPoints(IK, KJ) {

		IJ = bestChain(IJ, IK, KJ)
	}

	return IJ
}

func bestChain(IJ, IK, KJ chain) chain {
	if IJ.product < (IK.product * KJ.product) {
		IJ.product = IK.product * KJ.product
		IJ.path = joinPath(IK, KJ)
	}
	return IJ
}

//Проверяет две цепочки на отсутствие повторяющихся точек-валют в их "маршрутах"
func noMatchingPoints(IK, KJ chain) bool {

	visited := make(map[int]bool)

	for _, v := range IK.path {
		visited[v] = true
	}

	for _, v := range KJ.path {
		if visited[v] == true {
			return false
		}
	}

	return true
}

//Возвращает путь из I в J через K
func joinPath(IK, KJ chain) []int {

	IK.path = append(IK.path, KJ.path...)

	return IK.path
}
