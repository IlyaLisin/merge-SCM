package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"sort"

	"math/rand"
	"time"
	"math"
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

	byteString map[int]int
}

func main() {
	scm := new(SCM)
	scm.initSCM()
	scm.prepareRoutes()

	scm.generateInitState()

	initT := float64(10) // init T
	endT := 0.00001

	state := scm.byteString
	currentEnergy := calculateEnergy(state, scm.Routes)
	T := initT

	var stateCandidate map[int]int
	var candidateEnergy int
	var p float64

	for i := 0; i < 10000; i++ {
		stateCandidate = generateState(state)
		candidateEnergy = calculateEnergy(stateCandidate, scm.Routes)

		if candidateEnergy < currentEnergy {
			currentEnergy = candidateEnergy
			state = stateCandidate
		} else {
			p = transitionProbability(candidateEnergy - currentEnergy, T)

			if isTransit(p) {
				currentEnergy = candidateEnergy
				state = stateCandidate
			}
			}

			T = decreaseT(initT, i)

			if T <= endT {
				break
			}
	}

	fmt.Println("State: ", state)
	fmt.Println("Sum: ", calculateEnergy(state, scm.Routes))

	// time.Sleep(time.Duration(100)*time.Second)
}

func decreaseT(initT float64, i int) float64 {
	return initT * 0.1/float64(i)
}

func isTransit(p float64) bool {
	rand.Seed(time.Now().Unix())
	value := rand.Float64()

	if value <= p {
		return true
	}

	return false
}

func transitionProbability(deltaE int, T float64) float64 {
	return math.Exp(-float64((float64(deltaE)/T)))
}

func generateState(byteString map[int]int) map[int]int {
	rand.Seed(time.Now().Unix())

	i := rand.Intn(len(byteString))
	j := rand.Intn(len(byteString))

	if byteString[i] == 1 {
		byteString[i] = 0
	} else {
		byteString[i] = 1
	}

	if byteString[j] == 1 {
		byteString[j] = 0
	} else {
		byteString[j] = 1
	}

	return byteString
}

func calculateEnergy(byteString map[int]int, routes []Route) int {
	// This is Magic ╰( ⁰ ਊ ⁰ )━☆ﾟ.*･｡..::****::..
	// стваим только существующие роуты
	newRoutes := make([]Route, 0)
	for k := 0; k < len(byteString); k++ {
		if byteString[k] == 1 {
			for m, r := range routes {
				if r.ID == k {
					newRoutes = append(newRoutes, routes[m])
				}
			}
		}
	}

	sum := 0
	for _, r := range newRoutes {
		sum += r.S
	}

	return sum
}

func (scm *SCM) generateInitState() {
	scm.byteString = make(map[int]int, len(scm.Routes))

	rand.Seed(time.Now().Unix())

	// This is Magic ╰( ⁰ ਊ ⁰ )━☆ﾟ.*･｡..::****::..
	// стваим только существующие роуты
	newRoutes := make([]Route, 0)
	for i := 0; i < len(scm.Routes); i++ {
		scm.byteString[i] = rand.Intn(2)
		if scm.byteString[i] == 1 {
			for j, r := range scm.Routes {
				if r.ID == i {
					newRoutes = append(newRoutes, scm.Routes[j])
				}
			}
		}
	}
	scm.Routes = newRoutes
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
