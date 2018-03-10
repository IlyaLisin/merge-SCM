package ga

import (
	"service"
)

type Chromosome struct {
	gens map[int]bool
	fitness int
}

//var chromosome map[int]bool

func Init(gensCnt int) []*Chromosome{
	r := service.New()


	ch_arr := make([]*Chromosome,0)
	for n := 0; n < 5; n++ {
		ch := new(Chromosome)
		for i := 0; i < gensCnt; i++ {
			ch.gens[i] = r.Bool()
		}
		ch_arr = append(ch_arr, ch)
	}

	return ch_arr
}

func (chr *Chromosome) Fitness() {

}
