package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"sort"
	"math/rand"
	"time"
	"service"

	"github.com/ulule/deepcopier"
)

// V_COUNT - количество узлов, E_COUNT - количество ребер
const (
	V_COUNT = 9
	E_COUNT = 100

	P_COUNT = 3
	C_COUNT = 3
	S_COUNT = 3
)
//var E_COUNT int
//
//var P_COUNT int //производство
//var C_COUNT int //потребление
//var S_COUNT int //склад

// структура ребра
//type e_struct struct {
//	v1 v_struct
//	v2 v_struct
//	L  int
//	S  int
//	A  int
//
//	exists bool
//	Number int //номер на всем графе
//}
//
//// структура вершины
//type v_struct struct {
//	Number int
//	Type   int
//	Volume int
//	ID     int
//}
//
//// структура графа
//type SCM_struct struct {
//	v []v_struct
//	auto []e_struct
//	railway []e_struct
//	//available []e_struct
//
//	adj_matrix [81][81]e_struct
//
//	adj_matrix1 [P_COUNT][S_COUNT]e_struct
//	adj_matrix2 [S_COUNT][C_COUNT]e_struct
//
//	eCount int
//
//	chromosome *Chromosome
//}

type Chromosome struct {
	gens map[int]bool
	fitness float32
}

func main() {
	//scm := new(SCM)
	//scm.initSCM()
	//scm.generateRoutes()

	scm := new(SCM)
	scm.initSCM()
	scm.prepareRoutes()

	buf := *scm
	buf.Routes = nil

	//fmt.Println(scm.Routes)
	//fmt.Println(buf.Routes)

	population := make([]SCM, 5)
	r := service.New()

	//init 5 first
	for n := 0; n < 5; n++ {
		// TODO How copying struct 
		subScm := SCM{}
		deepcopier.Copy(scm).To(subScm)
		//subScm := *scm
		//subScm.Routes = append(scm.Routes)
		//subScm.Tops = append(scm.Tops)
		subScm.chromosome = new(Chromosome)
		subScm.chromosome.gens = make(map[int]bool, len(subScm.Routes))

		// This is Magic ╰( ⁰ ਊ ⁰ )━☆ﾟ.*･｡..::****::..
		// стваим только существующие роуты
		newRoutes := make([]Route, 0)
		for i := 0; i < len(subScm.Routes); i++ {
			subScm.chromosome.gens[i] = r.Bool()
			if subScm.chromosome.gens[i] {
				for j, r := range subScm.Routes {
					if r.ID == i {
						newRoutes = append(newRoutes, subScm.Routes[j])
					}
				}
			}
		}
		subScm.Routes = newRoutes

		for _, v := range subScm.Tops {
			v.NextIds = make([]int, 0)
			for _, r := range subScm.Routes {
				if v.ID == r.V1 {
					v.NextIds = append(v.NextIds, r.V2)
				}
			}
			fmt.Println(v.NextIds)
		}

		population = append(population, subScm)
	}

	//for _, s := range population {
	//	//fmt.Println(s.chromosome)
	//	fmt.Println(s.Routes)
	//	fmt.Println(s.Tops)
	//}

	//fmt.Println(population)
	fmt.Println("123123")
	fmt.Println(population[1].Routes)
	fmt.Println(population[1].Tops)

	//scm.initSCM()
	//
	//r := service.New()
	//
	//ch_arr := make([]*Chromosome,0)
	//for n := 0; n < 5; n++ {
	//	ch := new(Chromosome)
	//	ch.gens = make(map[int]bool,0)
	//	for i := 0; i < scm.eCount; i++ {
	//		ch.gens[i] = r.Bool()
	//	}
	//	ch_arr = append(ch_arr, ch)
	//}
	//
	//scm.chromosome = ch_arr[0]
	//
	//scm.build()
	////fmt.Println("QQQQ", scm.adj_matrix1)
	////os.Exit(1)
	//
	//for evol := 0; evol < 5; evol++ {
	//	fmt.Println("Evolution %d", evol)
	//	fmt.Println("ch_arr count ", len(ch_arr))
	//	for _, ch := range ch_arr {
	//		scm.chromosome = ch
	//		scm.build()
	//		ch.fitness = scm.Fitness()
	//		fmt.Println("Fitness", ch.fitness)
	//	}
	//
	//	//fmt.Println("%d", len(ch_arr))
	//
	//	// Селекция
	//	ch_arr = selection(ch_arr)
	//
	//	ch_arr = crossover(ch_arr)
	//
	//	//TODO мутации
	//}
	//
	//fmt.Println("Graph: ", scm.chromosome.gens)
	//fmt.Println("Graph: ", scm.adj_matrix1, " ", scm.adj_matrix2)
	//
	////
	//for i := range scm.adj_matrix[0] {
	//	for j := range scm.adj_matrix[0] {
	//		if scm.adj_matrix[i][j].S != 0 {
	//			fmt.Printf("%+v\n", scm.adj_matrix[i][j])
	//			fmt.Println("")
	//		}
	//	}
	//}

	//fmt.Println("двумерный массив: ", scm.adj_matrix)
}

type Vs struct {
	Uuid string `json:"uuid"`
	ID int `json:"id"`
	Type int `json:"type"`
	Volume int `json:"volume"`
	NextIds []int `json:"next_ids"`
}

type Route struct {
	ID int `json:"id"`
	V1 int `json:"v_id_1"`
	V2 int	`json:"v_id_2"`
	Type int	`json:"a"`
	L int `json:"l"`
	S int `json:"s"`
	Cells []int `json:"cells"`
}

type SCM struct {
	Tops []Vs `json:"tops"`
	Routes []Route	`json:"routes"`

	chromosome *Chromosome
}

// выбираем на каждый путь наименьшуюю стоимость
func (scm *SCM) prepareRoutes() {
	routes := readRoutes("src/config/new_routes.json")

	//fmt.Println(len(routes))

	minRoutes := make([]Route, 0)

	// This is Magic ╰( ⁰ ਊ ⁰ )━☆ﾟ.*･｡..::****::..
	excluded := make([]int, 0)
	for i, route := range routes {
		sameRoutes := make([]Route, 0)
		next := false
		for _, ex := range excluded {
			if i == ex {
				next = true
				break
			}
		}
		if next {
			continue
		}
		for j, subRoute := range routes {
			// находим одинаковые роуты
			if route.V1 == subRoute.V1 && route.V2 == subRoute.V2 {
				sameRoutes = append(sameRoutes, subRoute)
				excluded = append(excluded, j)
			}
		}
		minRoute := sameRoutes[0]
		if len(sameRoutes) > 1 {
			for _, r := range sameRoutes {
				if r.S < minRoute.S {
					minRoute = r
				}
			}
		}
		minRoutes = append(minRoutes, minRoute)
	}

	scm.Routes = minRoutes

	//countProd := 0
	//countStor := 0
	//countCons := 0
	//
	//for _, v := range scm.Tops {
	//	if v.Type == 0 {
	//		countProd += 1
	//	}
	//	if v.Type == 1 {
	//		countStor += 1
	//	}
	//	if v.Type == 2 {
	//		countCons += 1
	//	}
	//}

	// пронумеруем по порядку
	sort.Slice(scm.Routes, func(i, j int) bool {
		if scm.Routes[i].V1 < scm.Routes[j].V1 {
			return true
		}
		if scm.Routes[i].V1 > scm.Routes[j].V1 {
			return false
		}
		return scm.Routes[i].V2 < scm.Routes[j].V2
	})

	for i := range scm.Routes {
		scm.Routes[i].ID = i + 1
	}

}

func (scm *SCM) initSCM() *SCM {

	//for i := 1; i < 10; i++ {
	//	fmt.Println(uuid.V4())
	//}

	graph1 := readGraph("src/config/graph_1.json")
	graph2 := readGraph("src/config/graph_2.json")


	scm.Tops = append(graph1.Tops, graph2.Tops...)

	// сортировка по типу вершины
	sort.Slice(scm.Tops, func(i, j int) bool {
		return scm.Tops[i].Type < scm.Tops[j].Type
	})
	// записываем индексы по порядку
	for index := range scm.Tops {
		scm.Tops[index].ID = index + 1
	}

	//for _, element := range scm.Tops {
	//	println(element.ID)
	//}

	return scm
}

func readRoutes(path string) []Route {
	jsonFile, err := os.Open(path)

	//if err != nil {
	//	fmt.Println(err)
	//}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var routes []Route

	json.Unmarshal(byteValue, &routes)

	return routes
}

func readGraph(path string) SCM {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var scmJson SCM

	json.Unmarshal(byteValue, &scmJson)

	return scmJson
}

func (scm *SCM) generateRoutes() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	routes := []Route{}

	for i := 1; i < 101; i++ {
		for _, top := range scm.Tops {
			if (top.Type == 2) {
				continue
			}
			v2 := 0
			if (top.Type == 0) {
				v2 = r1.Intn(6 - 3) + 3
			}
			if (top.Type == 1) {
				v2 = r1.Intn(9 - 6) + 6
			}

			r := Route{
				ID:   i,
				Type: 0,
				V1: top.ID,
				V2: v2,
				L: r1.Intn(10),
				S: r1.Intn(10),
				Cells: make([]int, 0),
			}

			routes = append(routes, r)
		}
	}

	jsonData, err := json.Marshal(routes)
	//fmt.Println(routes)
	//fmt.Println(jsonData)
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("src/config/new_routes.json", jsonData, 0644)
}

func (scm *SCM) checkGraph() bool {



	return true
}
//
//func (scm *SCM_struct) checkGraph_2() bool {
//	tmpV := make(map[int]int,0)
//
//	for i := range scm.adj_matrix1[0] {
//		for j := range scm.adj_matrix1[0] {
//			// проверка на производство -> склад -> потребление
//			if scm.adj_matrix1[i][j].exists && scm.adj_matrix1[i][j].v1.Type >= scm.adj_matrix1[i][j].v2.Type {
//				return false
//			}
//			tmpV[scm.adj_matrix1[i][j].v1.ID] = 1
//			tmpV[scm.adj_matrix1[i][j].v2.ID] = 1
//		}
//	}
//	for i := range scm.adj_matrix2[0] {
//		for j := range scm.adj_matrix2[0] {
//			// проверка на производство -> склад -> потребление
//			if scm.adj_matrix2[i][j].exists && scm.adj_matrix2[i][j].v1.Type >= scm.adj_matrix2[i][j].v2.Type {
//				return false
//			}
//			tmpV[scm.adj_matrix2[i][j].v1.ID] = 1
//			tmpV[scm.adj_matrix2[i][j].v2.ID] = 1
//		}
//	}
//
//	// Все ли вершины задействованы
//	var b bool
//	for i := 0; i < V_COUNT; i++ {
//		b = false
//		for id, exist := range tmpV {
//			if i == id && exist == 1 {
//				b = true
//			}
//		}
//		if !b {
//			return false
//		}
//	}
//
//	return true
//}
//
//func (scm *SCM_struct) Fitness() float32 {
//	// берем хромосому и если ген=0, ставим туда пустой роут с exist: false
//	for number, exist := range scm.chromosome.gens {
//		if !exist {
//			all_routes := append(scm.auto, scm.railway...)
//			for i := 0; i < len(all_routes) ;i++ {
//				if all_routes[i].Number == number {
//					if all_routes[i].v1.Type == 0 {
//						scm.adj_matrix1[all_routes[i].v1.Number][all_routes[i].v2.Number] = e_struct{exists:false}
//					}
//					if all_routes[i].v1.Type == 1 {
//						scm.adj_matrix2[all_routes[i].v1.Number][all_routes[i].v2.Number] = e_struct{exists:false}
//					}
//				}
//			}
//		}
//	}
//
//	if !scm.checkGraph_2() {
//
//		return 0
//	}
//
//	total := 0
//	for i := 0; i < len(scm.adj_matrix1[0]); i++ {
//		for j := 0; j < len(scm.adj_matrix1[0]); j++ {
//			total += scm.adj_matrix1[i][j].S
//		}
//	}
//	for i := 0; i < len(scm.adj_matrix2[0]); i++ {
//		for j := 0; j < len(scm.adj_matrix2[0]); j++ {
//			total += scm.adj_matrix2[i][j].S
//		}
//	}
//
//	return 1/float32(total)
//}

func selection(chromosomes []*Chromosome) []*Chromosome{
	//fmt.Println("selection")
	var totalFitness float32 = 0
	for _, ch := range chromosomes {
		totalFitness += ch.fitness
	}

	averageFitness := totalFitness/float32(len(chromosomes))
	// Выбираем больше ср значения
	var ret []*Chromosome
	for _, ch := range chromosomes {
		if ch.fitness >= averageFitness {
			ret = append(ret, ch)
		}
	}
	if len(ret) < 5 {
		// Дописываем до количества
		for ; len(ret) < len(chromosomes); {
			for _, ch := range chromosomes {
				if ch.fitness >= averageFitness {
					ret = append(ret, ch)
				}
			}
		}
	}

	return ret
}

func crossover(chromosomes []*Chromosome) []*Chromosome {
	//fmt.Println("crossover")
	childs := make([]*Chromosome,0)
	i := 0
	j := 1
	//fmt.Println("1 ", chromosomes[0].gens)
	//fmt.Println("1 ", chromosomes[1].gens)
	//fmt.Println("1 ", chromosomes[2].gens)
	//fmt.Println("1 ", chromosomes[3].gens)
	//fmt.Println("1 ", chromosomes[4].gens)
	for ; j < len(chromosomes) ; {
		p := len(chromosomes[0].gens)/2
		tmpCh1 := new(Chromosome)
		tmpCh1.gens = make(map[int]bool,0)
		tmpCh2 := new(Chromosome)
		tmpCh2.gens = make(map[int]bool,0)
		for m := 0; m < len(chromosomes[0].gens); m++ {
			if m <= p {
				tmpCh1.gens[m] = chromosomes[i].gens[m]
				tmpCh2.gens[m] = chromosomes[j].gens[m]
			} else {
				tmpCh1.gens[m] = chromosomes[j].gens[m]
				tmpCh2.gens[m] = chromosomes[i].gens[m]
			}
		}
		childs = append(childs, tmpCh1, tmpCh2)
		i += 2
		j += 2
	}

	if len(childs) < 5 {
		if len(childs) != len(chromosomes) {
			childs = append(childs, chromosomes[0])
		}
	}

	//fmt.Println("2 ", childs[0].gens)
	//fmt.Println("2 ", childs[1].gens)
	//fmt.Println("2 ", childs[2].gens)
	//fmt.Println("2 ", childs[3].gens)
	//fmt.Println("2 ", childs[4].gens)
	return childs
}
