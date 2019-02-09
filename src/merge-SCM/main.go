package main

import (
	"fmt"

	"math/rand"
	"time"
	"service"
)

// V_COUNT - количество узлов, E_COUNT - количество ребер
const (
	V_COUNT = 81
	E_COUNT = 100

	P_COUNT = 27
	C_COUNT = 27
	S_COUNT = 27
)
//var E_COUNT int
//
//var P_COUNT int //производство
//var C_COUNT int //потребление
//var S_COUNT int //склад

// структура ребра
type e_struct struct {
	v1 v_struct
	v2 v_struct
	L  int
	S  int
	A  int

	exists bool
	Number int //номер на всем графе
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

	adj_matrix1 [P_COUNT][S_COUNT]e_struct
	adj_matrix2 [S_COUNT][C_COUNT]e_struct

	eCount int

	chromosome *Chromosome
}

type Chromosome struct {
	gens map[int]bool
	fitness float32
}

func main() {
	t0 := time.Now()

	scm := new(SCM_struct)
	scm.initSCM_2()

	r := service.New()

	ch_arr := make([]*Chromosome,0)
	for n := 0; n < 5; n++ {
		ch := new(Chromosome)
		ch.gens = make(map[int]bool,0)
		for i := 0; i < scm.eCount; i++ {
			ch.gens[i] = r.Bool()
		}
		ch_arr = append(ch_arr, ch)
	}

	scm.chromosome = ch_arr[0]

	scm.build()
	//fmt.Println("QQQQ", scm.adj_matrix1)
	//os.Exit(1)

	for evol := 0; evol < 5; evol++ {
		fmt.Println("Evolution %d", evol)
		fmt.Println("ch_arr count ", len(ch_arr))
		for _, ch := range ch_arr {
			scm.chromosome = ch
			scm.build()
			ch.fitness = scm.Fitness()
			fmt.Println("Fitness", ch.fitness)
		}

		//fmt.Println("%d", len(ch_arr))

		// Селекция
		ch_arr = selection(ch_arr)

		ch_arr = crossover(ch_arr)

		//TODO мутации
	}

	fmt.Println("Graph: ", scm.chromosome.gens)
	fmt.Println("Graph: ", scm.adj_matrix1, " ", scm.adj_matrix2)

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

	t1 := time.Now()
	fmt.Printf("Elapsed time: %v", t1.Sub(t0))
}

//func (scm *SCM_struct) initSCM() *SCM_struct {
//	s1 := rand.NewSource(time.Now().UnixNano())
//	r1 := rand.New(s1)
//
//	//scm := new(SCM_struct)
//	fmt.Println("init V")
//	// генерация узлов, Vij
//	for i := 0; i < V_COUNT; i++ {
//		v := v_struct{
//			Number: i%3,
//			Type: i/3,
//			ID: i,
//		}
//		scm.v = append(scm.v, v)
//		fmt.Println(scm.v[i].Number, " ", scm.v[i].Type)
//	}
//
//	eNumber := -1
//
//	// генерация ребер автотранспорта
//	for i := 0; i < E_COUNT; i++ {
//		eNumber += 1
//		e := e_struct{
//			v1: scm.v[r1.Intn(V_COUNT)],
//			v2: scm.v[r1.Intn(V_COUNT)],
//			L: r1.Intn(10),
//			S: r1.Intn(10),
//			A: 0,
//			exists: true,
//			Number: eNumber,
//		}
//		scm.auto = append(scm.auto, e)
//	}
//
//	// генерация ребер ж/д транспорта
//	for i := 0; i < E_COUNT; i++ {
//		eNumber += 1
//		e := e_struct{
//			v1: scm.v[r1.Intn(V_COUNT)],
//			v2: scm.v[r1.Intn(V_COUNT)],
//			L: r1.Intn(10),
//			S: r1.Intn(10),
//			A: 1,
//			exists: true,
//			Number: eNumber,
//		}
//		scm.railway = append(scm.railway, e)
//	}
//
//	//for i := 0; i < 10; i++ {
//	//	e := e_struct{
//	//		v1: r1.Intn(10),
//	//		v2: r1.Intn(10),
//	//		L: r1.Intn(10),
//	//		S: r1.Intn(10),
//	//	}
//	//	scm.available[i] = e
//	//}
//
//	//for i := range scm.v {
//	//	for j := range scm.v {
//	//		scm.adj_matrix = append(scm.adj_matrix[i], e_struct{})
//	//	}
//	//}
//
//	// все роуты
//	all_routes := append(scm.auto, scm.railway...)
//	// выбор роутов с наим стоимостью
//	for i := range all_routes {
//		v1 := all_routes[i].v1
//		v2 := all_routes[i].v2
//		if v1 == v2 { continue }
//		// если 0,0 или стоимость ниже уже записанной
//		if (((scm.adj_matrix[v1.ID][v2.ID].v1 == v_struct{} && scm.adj_matrix[v1.ID][v2.ID].v2 == v_struct{}) ||
//			(scm.adj_matrix[v1.ID][v2.ID].S > all_routes[i].S)) &&
//			((v1.Number != v2.Number)) &&
//			(v2.Type - v1.Type == 1)) {
//			scm.adj_matrix[v1.ID][v2.ID] = all_routes[i]
//		}
//	}
//
//	scm.eCount = V_COUNT*2
//
//	//for i := range scm.adj_matrix[0] {
//	//	for j := range scm.adj_matrix[0] {
//	//		// Incostiling
//	//		if scm.adj_matrix[i][j].S != 0 && scm.adj_matrix[j][i].S != 0 && scm.adj_matrix[i][j].S >= scm.adj_matrix[j][i].S {
//	//			scm.adj_matrix[i][j] = e_struct{}
//	//		}
//	//	}
//	//}
//
//	return scm
//}

func (scm *SCM_struct) initSCM_2() *SCM_struct {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	//n := P_COUNT + S_COUNT +  C_COUNT
	//pcn := P_COUNT + C_COUNT
	//sn := S_COUNT
	// число возможных ребер
	//r_count := (n*(n-1)/2) - ((pcn)*(pcn-1)/2) - (sn*(sn-1)/2)

	var pArr []v_struct
	var cArr []v_struct
	var sArr []v_struct

	//scm := new(SCM_struct)
	fmt.Println("init V")
	// генерация узлов, Vij
	for i := 0; i < V_COUNT; i++ {
		v := v_struct{
			Number: i%3,
			Type: i/3,
			ID: i,
		}
		switch v.Type {
		case 0:
			pArr = append(pArr, v)
		case 1:
			sArr = append(sArr, v)
		case 2:
			cArr = append(cArr, v)
		}
		scm.v = append(scm.v, v)
		fmt.Println(scm.v[i].Number, " ", scm.v[i].Type)
	}

	// генерация ребер автотранспорта
	for j := 0; j < E_COUNT; j++ {
		//fmt.Println("RANDOM ", r1.Intn(V_COUNT))
		e := e_struct{
			v1: pArr[r1.Intn(V_COUNT/3)],
			v2: sArr[r1.Intn(V_COUNT/3)],
			L:      r1.Intn(10),
			S:      r1.Intn(10),
			A:      0,
			exists: true,
		}
		e.Number = e.v1.Number*S_COUNT + e.v2.Number
		scm.auto = append(scm.auto, e)

		e = e_struct{
			v1: sArr[r1.Intn(V_COUNT/3)],
			v2: cArr[r1.Intn(V_COUNT/3)],
			L:      r1.Intn(10),
			S:      r1.Intn(10),
			A:      0,
			exists: true,
		}
		e.Number = e.v1.Number*S_COUNT + e.v2.Number + V_COUNT
		scm.auto = append(scm.auto, e)
	}

	//// генерация ребер ж/д транспорта
	//for i := 0; i < E_COUNT; i++ {
	//	eNumber += 1
	//	e := e_struct{
	//		v1: scm.v[r1.Intn(V_COUNT)],
	//		v2: scm.v[r1.Intn(V_COUNT)],
	//		L: r1.Intn(10),
	//		S: r1.Intn(10),
	//		A: 1,
	//		exists: true,
	//	}
	//	switch e.v1.Type {
	//	case 0:
	//		e.Number = e.v1.Number*S_COUNT + e.v2.Number
	//	case 1:
	//		e.Number = e.v1.Number*S_COUNT + e.v2.Number + V_COUNT
	//	default:
	//		e.Number = -1
	//	}
	//	scm.railway = append(scm.railway, e)
	//}

	scm.eCount = V_COUNT*2

	return scm
}

func (scm *SCM_struct) build() {
	all_routes := append(scm.auto, scm.railway...)

	for i := 0; i < len(scm.adj_matrix1[0]); i++ {
		for j := 0; j < len(scm.adj_matrix1[0]); j++ {
			scm.adj_matrix1[i][j] = e_struct{exists:false}
		}
	}
	for i := 0; i < len(scm.adj_matrix2[0]); i++ {
		for j := 0; j < len(scm.adj_matrix2[0]); j++ {
			scm.adj_matrix2[i][j] = e_struct{exists:false}
		}
	}

	for n, exist := range scm.chromosome.gens {
		if exist {
			for _, r := range all_routes {
				// Ищет роут с таким номером ребра и пишем
				if r.Number == n {
					if r.v1.Type == 0 {
						scm.adj_matrix1[r.v1.Number][r.v2.Number] = r
					}
					if r.v1.Type == 1 {
						scm.adj_matrix2[r.v1.Number][r.v2.Number] = r
					}
					break
				}
			}
		}
	}
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

func (scm *SCM_struct) checkGraph_2() bool {
	tmpV := make(map[int]int,0)

	for i := range scm.adj_matrix1[0] {
		for j := range scm.adj_matrix1[0] {
			// проверка на производство -> склад -> потребление
			if scm.adj_matrix1[i][j].exists && scm.adj_matrix1[i][j].v1.Type >= scm.adj_matrix1[i][j].v2.Type {
				return false
			}
			tmpV[scm.adj_matrix1[i][j].v1.ID] = 1
			tmpV[scm.adj_matrix1[i][j].v2.ID] = 1
		}
	}
	for i := range scm.adj_matrix2[0] {
		for j := range scm.adj_matrix2[0] {
			// проверка на производство -> склад -> потребление
			if scm.adj_matrix2[i][j].exists && scm.adj_matrix2[i][j].v1.Type >= scm.adj_matrix2[i][j].v2.Type {
				return false
			}
			tmpV[scm.adj_matrix2[i][j].v1.ID] = 1
			tmpV[scm.adj_matrix2[i][j].v2.ID] = 1
		}
	}

	// Все ли вершины задействованы
	var b bool
	for i := 0; i < V_COUNT; i++ {
		b = false
		for id, exist := range tmpV {
			if i == id && exist == 1 {
				b = true
			}
		}
		if !b {
			return false
		}
	}

	return true
}

func (scm *SCM_struct) Fitness() float32 {
	// берем хромосому и если ген=0, ставим туда пустой роут с exist: false
	for number, exist := range scm.chromosome.gens {
		if !exist {
			all_routes := append(scm.auto, scm.railway...)
			for i := 0; i < len(all_routes) ;i++ {
				if all_routes[i].Number == number {
					if all_routes[i].v1.Type == 0 {
						scm.adj_matrix1[all_routes[i].v1.Number][all_routes[i].v2.Number] = e_struct{exists:false}
					}
					if all_routes[i].v1.Type == 1 {
						scm.adj_matrix2[all_routes[i].v1.Number][all_routes[i].v2.Number] = e_struct{exists:false}
					}
				}
			}
		}
	}

	if !scm.checkGraph_2() {

		return 0
	}

	total := 0
	for i := 0; i < len(scm.adj_matrix1[0]); i++ {
		for j := 0; j < len(scm.adj_matrix1[0]); j++ {
			total += scm.adj_matrix1[i][j].S
		}
	}
	for i := 0; i < len(scm.adj_matrix2[0]); i++ {
		for j := 0; j < len(scm.adj_matrix2[0]); j++ {
			total += scm.adj_matrix2[i][j].S
		}
	}

	return 1/float32(total)
}

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
