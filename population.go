package goevo

import (
	"math/rand"
	"sort"
)

type Agent struct {
	GT      *Genotype
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

func (p Population) Repopulate(ratio float64, f func(g1 *Genotype, g2 *Genotype) *Genotype) {
	p.Order()
	numberToKeep := int(ratio * float64(len(p)))
	for g := numberToKeep; g < len(p); g++ {
		parent1I := rand.Intn(numberToKeep)
		parent2I := rand.Intn(numberToKeep)
		gt := f(p[parent1I].GT, p[parent2I].GT)
		pt := GrowPhenotype(gt)
		p[g] = &Agent{GT: gt, PT: pt}
	}
}

func CreatePopulation(gts []*Genotype) Population {
	p := make(Population, len(gts))
	for i := range p {
		pt := GrowPhenotype(gts[i])
		p[i] = &Agent{GT: gts[i], PT: pt}
	}
	return p
}
