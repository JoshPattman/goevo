package goevo

import (
	"math/rand"
	"sort"
)

type Agent struct {
	GT      *GenotypeSlow
	PT      *Phenotype
	Fitness float64
}

type Population []*Agent

func (p Population) Order() {
	sort.Slice(p, func(i, j int) bool {
		// sort by fittest first
		return p[i].Fitness > p[j].Fitness
	})
}

func (p Population) Repopulate(ratio float64, f func(g1 *GenotypeSlow, g2 *GenotypeSlow) *GenotypeSlow) {
	p.Order()
	numberToKeep := int(ratio * float64(len(p)))
	for g := numberToKeep; g < len(p); g++ {
		parent1I := rand.Intn(numberToKeep)
		parent2I := rand.Intn(numberToKeep)
		if p[parent1I].Fitness < p[parent2I].Fitness {
			c := parent2I
			parent2I = parent1I
			parent1I = c
		}
		gt := f(p[parent1I].GT, p[parent2I].GT)
		pt := GrowPhenotypeLegacy(gt)
		p[g] = &Agent{GT: gt, PT: pt}
	}
}

func NewPopulation(gts []*GenotypeSlow) Population {
	p := make(Population, len(gts))
	for i := range p {
		pt := GrowPhenotypeLegacy(gts[i])
		p[i] = &Agent{GT: gts[i], PT: pt}
	}
	return p
}
