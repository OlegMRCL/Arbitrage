package exchange


//Данные о бирже
type Exchange struct {
	Currencies []string		//список валют на бирже
	PriceTable PriceTable	//таблица межвалютных цен
	Fee float64				//комиссия, взимаемая на бирже
}


//Список валютных пар
type PairList map[string]Pair


//Данные о валютной паре
type Pair struct {
	Bid float64		//цена спроса
	Ask float64		//цена предложения
}


//В таблице PriceTable храянятся цены валют:
// в ячейке [i][j] указана цена валюты i в валюте j,
// а в ячейке [j][i] - валюты j в валюте i.
type PriceTable [][] float64


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
func (e *Exchange) FindArbitrage() (results Results){

	at := e.generateTable(e.PriceTable)

	numCurr := len(e.PriceTable)

	for k := 0; k < numCurr; k++ {
		for i := 0; i < numCurr; i++ {
			for j := 0; j < numCurr; j++ {

				if (i==j) && (at[i][j].product > 1) {

					path := e.Currencies[i]
					point := at[i][j].nextLink

					for point != i {
						path = "-->" + e.Currencies[point]
						point = at[point][j].nextLink
					}

					results[path] = at[i][j].product
				}

				at[i][j] = at.checkChains(i, j, k)
			}
		}
	}

	return results
}


type Results map[string]float64

//Данные о цепочке, признанной выгодной (Profit > 1)
type Result struct {
	Profit float64		//относительная величина профита
	Path []int			//валютный "маршрут"
}


//Возвращает цепочку, выбранную как наивыгоднейшую:
// выбирает между исходной цепочкой IJ и предлагаемой IKJ, сочлененной из IK и KJ.
func (at *ArbitrageTable) checkChains(i, j, k int) chain  {
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
		IJ.nextLink = IK.nextLink
		IJ.visited = joinVisited(IK, KJ)
	}
	return IJ
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







