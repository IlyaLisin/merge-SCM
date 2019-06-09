package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"sort"
	"math"
	"strconv"
)

// V_COUNT - количество узлов, E_COUNT - количество ребер
const (
	P_COUNT = 3
	C_COUNT = 3
	S_COUNT = 3
)

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
}

func main() {
	scm := new(SCM)
	scm.initSCM()
	scm.prepareRoutes()

	max := int(math.Pow(2, 18))

	byteString := make(map[int]int, len(scm.Routes))

	for i := 0; i < len(scm.Routes); i++ {
		byteString[i] = 0
	}

	bestSum := 9999999
	bestByteString := make(map[int]int, len(scm.Routes))

	for i := 0; i < max; i++ {
		str := strconv.FormatInt(int64(i), 2)

		for k, s := range str {
			if s == '1' {
				byteString[k] = 1
			}
		}

		// This is Magic ╰( ⁰ ਊ ⁰ )━☆ﾟ.*･｡..::****::..
		// стваим только существующие роуты
		newRoutes := make([]Route, 0)
		for k := 0; k < len(scm.Routes); k++ {
			if byteString[k] == 1 {
				for m, r := range scm.Routes {
					if r.ID == k {
						newRoutes = append(newRoutes, scm.Routes[m])
					}
				}
			}
		}

		if !isAllowed(newRoutes) {
			continue
		}

		sum := 0
		for k := 0; k < len(newRoutes); k++ {
			sum += newRoutes[k].S
		}

		if sum < bestSum {
			bestSum = sum
			bestByteString = byteString
			fmt.Println("Best", bestByteString)
		}

		for i := 0; i < len(scm.Routes); i++ {
			byteString[i] = 0
		}
	}

	//fmt.Println("Best", bestByteString)
	fmt.Println("Best", bestSum)

	// time.Sleep(time.Duration(100)*time.Second)
}

func isAllowed(routes []Route) bool {
	// Отсекаем графы, в которых не все вершины задействованы
	b := false
	var excluded bool

	for i := 0; i < P_COUNT + S_COUNT; i++  {
		b = false
		for _, r := range routes {
			if r.V1 == i {
				b = true
				break
			}
		}
		if b {
			continue
		}

		excluded = true
	}

	if excluded {
		return false
	}

	for i := P_COUNT + S_COUNT; i < P_COUNT + C_COUNT + S_COUNT; i++  {
		b = false
		for _, r := range routes {
			if r.V2 == i {
				b = true
				break
			}
		}
		if b {
			continue
		}

		excluded = true
	}

	if excluded {
		return false
	}

	return true
}

// выбираем на каждый путь наименьшуюю стоимость
func (scm *SCM) prepareRoutes() {
	routes := readRoutes("src/config/routes.json")

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
	graph := readGraph("src/config/graph.json")

	scm.Tops = graph.Tops

	// сортировка по типу вершины
	sort.Slice(scm.Tops, func(i, j int) bool {
		return scm.Tops[i].Type < scm.Tops[j].Type
	})
	// записываем индексы по порядку
	for index := range scm.Tops {
		scm.Tops[index].ID = index + 1
	}

	return scm
}

func readRoutes(path string) []Route {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

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
