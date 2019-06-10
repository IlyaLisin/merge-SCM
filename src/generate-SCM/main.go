package main

import (
	"math/rand"
	"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// V_COUNT - количество узлов, E_COUNT - количество ребер
const (
	V_COUNT = 50

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
}

type SCM struct {
	Tops []Vs `json:"tops"`
	Routes []Route	`json:"routes"`
}

func main() {
	scm := new(SCM)
	scm.generateGraph()
	scm.generateRoutes()
}

func (scm *SCM) generateRoutes() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	routes := []Route{}

	for i := 1; i < V_COUNT*10; i++ {
		for _, top := range scm.Tops {
			if (top.Type == 2) {
				continue
			}
			v2 := 0
			if (top.Type == 0) {
				v2 = r1.Intn(S_COUNT) + P_COUNT
			}
			if (top.Type == 1) {
				v2 = r1.Intn(C_COUNT) + P_COUNT + S_COUNT
			}

			if top.ID >= v2 {
				continue
			}

			r := Route{
				ID:   i,
				Type: 0,
				V1: top.ID,
				V2: v2,
				L: r1.Intn(10),
				S: 10,
			}

			routes = append(routes, r)
		}
	}

	jsonData, err := json.Marshal(routes)

	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("src/config/routes50.json", jsonData, 0644)

	if err != nil {
		fmt.Println(err)
	}
}

func (scm *SCM) generateGraph() {
	tops := []Vs{}

	j := 0
	for i := 0; i < P_COUNT; i++ {
		top :=  Vs {
			ID: j,
			Type: 0,
			Volume: 50,
		}
		tops = append(tops, top)
		j++
	}

	for i := 0; i < S_COUNT; i++ {
		top :=  Vs {
			ID: j,
			Type: 1,
			Volume: 50,
		}
		tops = append(tops, top)
		j++
	}

	for i := 0; i < C_COUNT; i++ {
		top :=  Vs {
			ID: j,
			Type: 2,
			Volume: 50,
		}
		j++
		tops = append(tops, top)
	}

	scm.Tops = tops

	jsonData, err := json.Marshal(tops)

	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("src/config/graph50.json", jsonData, 0644)

	if err != nil {
		fmt.Println(err)
	}
}
