package ga

import (
	"service"
)



type Chromosome struct {
	gens map[int]bool
	fitness int
}

type GA struct {
	Chromosomes []*Chromosome
}

//var chromosome map[int]bool

func (ga *GA) Init(chrCnt int, allelsCnt int) {
	r := service.New()


	ch_arr := make([]*Chromosome,0)
	for n := 0; n < chrCnt; n++ {
		ch := new(Chromosome)
		for i := 0; i < allelsCnt; i++ {
			ch.gens[i] = r.Bool()
		}
		ch_arr = append(ch_arr, ch)
	}

	ga.Chromosomes = ch_arr
}
