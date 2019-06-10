package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"sort"
	"service"

	"github.com/jinzhu/copier"
	"time"
)

// V_COUNT - количество узлов, E_COUNT - количество ребер
const (
	P_COUNT = 15
	C_COUNT = 20
	S_COUNT = 15
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

	chromosome *Chromosome
}

type Chromosome struct {
	gens map[int]bool
	fitness int
}

func main() {
	t0 := time.Now()
	scm := new(SCM)
	scm.initSCM()
	scm.prepareRoutes()

	buf := *scm
	buf.Routes = nil

	population := make([]*SCM, 0)
	r := service.New()

	//init 5 first
	for n := 0; n < 5; n++ {
		// TODO How copying struct
		subScm := SCM{}
		copier.Copy(&subScm.Tops, &scm.Tops)
		copier.Copy(&subScm.Routes, &scm.Routes)

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
			fmt.Print(v.NextIds)
		}
		fmt.Println("")

		population = append(population, &subScm)
	}

	chromosomes := make([]Chromosome, 0)

	for _, subScm := range population {
		subScm.fitness()
		chromosomes = append(chromosomes, *subScm.chromosome)
	}

	for evol := 0; evol < 1000; evol++ {
		//fmt.Println("Evolution", evol)

		childs := crossover(chromosomes)

		new_population := make([]*SCM, len(childs))

		for n, c := range childs {
			// This is Magic ╰( ⁰ ਊ ⁰ )━☆ﾟ.*･｡..::****::..
			// стваим только существующие роуты
			newRoutes := make([]Route, 0)
			for i := 0; i < len(scm.Routes); i++ {
				if c.gens[i] {
					for j, r := range scm.Routes {
						if r.ID == i {
							newRoutes = append(newRoutes, scm.Routes[j])
						}
					}
				}
			}

			new_population[n] = new(SCM)
			new_population[n].chromosome = c
			new_population[n].Routes = newRoutes
			new_population[n].fitness()
		}

		new_population = selection(new_population)

		chromosomes = make([]Chromosome, 0)

		for _, p := range new_population {
			chromosomes = append(chromosomes, *p.chromosome)
		}

		//TODO мутации
	}

	fmt.Println("Graph: ", scm.chromosome.gens)

	fmt.Println("\n")

	for _, s := range population {
		fmt.Println(s.Routes)
		fmt.Println(s.chromosome.fitness)
		fmt.Println("\n")
	}

	t1 := time.Now();

	fmt.Printf("Elapsed time: %v \n", t1.Sub(t0))
	// time.Sleep(time.Duration(100)*time.Second)
}

func (scm *SCM) fitness() {
	sum := 0
	for _, r := range scm.Routes {
		sum += r.S
	}

	scm.chromosome.fitness = sum
}

// выбираем на каждый путь наименьшуюю стоимость
func (scm *SCM) prepareRoutes() {
	routes := readRoutes("src/config/routes50.json")

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
	graph := readGraph("src/config/graph50.json")

	scm.Tops = graph.Tops

	// сортировка по типу вершины
	sort.Slice(scm.Tops, func(i, j int) bool {
		return scm.Tops[i].Type < scm.Tops[j].Type
	})
	// записываем индексы по порядку
	for index := range scm.Tops {
		scm.Tops[index].ID = index + 1
	}

	scm.chromosome = new(Chromosome)

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

func selection(scms []*SCM) []*SCM {
	matched := make([]*SCM, 0)
	// Отсекаем графы, в которых не все вершины задействованы
	for _, s := range scms {

		b := false
		var excluded bool

		for i := 0; i < P_COUNT + S_COUNT; i++  {
			b = false
			for _, r := range s.Routes {
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
			continue
		}

		for i := P_COUNT + S_COUNT; i < P_COUNT + C_COUNT + S_COUNT; i++  {
			b = false
			for _, r := range s.Routes {
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
			continue
		}

		matched = append(matched, s)

	}

	sort.Slice(matched, func(i, j int) bool {
		return matched[i].chromosome.fitness < matched[j].chromosome.fitness
	})

	if len(matched) >= 5 {
		return matched[0:4]
	} else {
		return matched
	}
}

func crossover(chromosomes []Chromosome) []*Chromosome {
	childs := make([]*Chromosome,0)

	i := 0
	j := 1

	for i = 0; i < len(chromosomes); i++ {
		for j = 1; j < len(chromosomes); j++ {
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
		}
	}

	return childs
}
