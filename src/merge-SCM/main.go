package main

import (
	"fmt"

	"math/rand"
	"time"
	"ga"
)

// V_COUNT - количество узлов, E_COUNT - количество ребер
var V_COUNT int
var E_COUNT int

// структура ребра
type e_struct struct {
	v1 v_struct
	v2 v_struct
	L  int
	S  int
	A  int
}

// структура вершины
type v_struct struct {
	Number int
	Type   int
	Volume int
	ID     int
}

// структура графа
type SCM_struct struct {
	v []v_struct
	auto []e_struct
	railway []e_struct
	//available []e_struct

	adj_matrix [81][81]e_struct

	vCount int
}

type GA struct {
	ga.Chromosome
}

func main() {
	V_COUNT = 9
	E_COUNT = 100

	scm := new(SCM_struct)
	scm.initSCM()

	chrms := ga.Init(5)

	//
	for i := range scm.adj_matrix[0] {
		for j := range scm.adj_matrix[0] {
			if scm.adj_matrix[i][j].S != 0 {
				fmt.Printf("%+v\n", scm.adj_matrix[i][j])
				fmt.Println("")
			}
		}
	}

	//fmt.Println("двумерный массив: ", scm.adj_matrix)
}

func (scm *SCM_struct) initSCM() *SCM_struct {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	//scm := new(SCM_struct)
	fmt.Println("init V")
	// генерация узлов, Vij
	for i := 0; i < V_COUNT; i++ {
		v := v_struct{
			Number: i%3,
			Type: i/3,
			ID: i,
		}
		scm.v = append(scm.v, v)
		fmt.Println(scm.v[i].Number, " ", scm.v[i].Type)
	}

	// генерация ребер автотранспорта
	for i := 0; i < E_COUNT; i++ {
		e := e_struct{
			v1: scm.v[r1.Intn(V_COUNT)],
			v2: scm.v[r1.Intn(V_COUNT)],
			L: r1.Intn(10),
			S: r1.Intn(10),
			A: 0,
		}
		scm.auto = append(scm.auto, e)
	}

	// генерация ребер ж/д транспорта
	for i := 0; i < E_COUNT; i++ {
		e := e_struct{
			v1: scm.v[r1.Intn(V_COUNT)],
			v2: scm.v[r1.Intn(V_COUNT)],
			L: r1.Intn(10),
			S: r1.Intn(10),
			A: 1,
		}
		scm.railway = append(scm.railway, e)
	}

	//for i := 0; i < 10; i++ {
	//	e := e_struct{
	//		v1: r1.Intn(10),
	//		v2: r1.Intn(10),
	//		L: r1.Intn(10),
	//		S: r1.Intn(10),
	//	}
	//	scm.available[i] = e
	//}

	//for i := range scm.v {
	//	for j := range scm.v {
	//		scm.adj_matrix = append(scm.adj_matrix[i], e_struct{})
	//	}
	//}

	// все роуты
	all_routes := append(scm.auto, scm.railway...)
    // выбор роутов с наим стоимостью
	for i := range all_routes {
		v1 := all_routes[i].v1
		v2 := all_routes[i].v2
		if v1 == v2 { continue }
		// если 0,0 или стоимость ниже уже записанной
		if (((scm.adj_matrix[v1.ID][v2.ID].v1 == v_struct{} && scm.adj_matrix[v1.ID][v2.ID].v2 == v_struct{}) ||
			(scm.adj_matrix[v1.ID][v2.ID].S > all_routes[i].S)) &&
			((v1.Number != v2.Number)) &&
			(v2.Type - v1.Type == 1)) {
			scm.adj_matrix[v1.ID][v2.ID] = all_routes[i]
		}
	}

	//for i := range scm.adj_matrix[0] {
	//	for j := range scm.adj_matrix[0] {
	//		// Incostiling
	//		if scm.adj_matrix[i][j].S != 0 && scm.adj_matrix[j][i].S != 0 && scm.adj_matrix[i][j].S >= scm.adj_matrix[j][i].S {
	//			scm.adj_matrix[i][j] = e_struct{}
	//		}
	//	}
	//}

	return scm
}

func (scm *SCM_struct) checkGraph() bool {
	// проверка на саму себя
	for i := range scm.adj_matrix[0] {
		if scm.adj_matrix[i][i].S != 0  {
			return false
		}
		for j := range scm.adj_matrix[0] {
			if i == j {
				continue
			}
			// проверка на производство -> склад -> потребление
			if scm.adj_matrix[i][j].v1.Type >= scm.adj_matrix[i][j].v2.Type {
				return false
			}
		}
	}

	return true
}

func (chr *GA) Fitness() int {

}

//func checkConnectivity(scm *SCM_struct) bool {
//	for
//}