package goevo

import (
	"math/rand"
	"sort"
)

type Agent struct {
	Genotype  Genotype
	Phenotype *Phenotype
	Fitness   float64
}

type Population []*Agent

func (p Population) Order() {
	sort.Slice(p, func(i, j int) bool {
		// sort by fittest first
		return p[i].Fitness > p[j].Fitness
	})
}

func (p Population) Repopulate(ratio float64, f func(g1 Genotype, g2 Genotype) Genotype, info ActivationInfo) {
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
		gt := f(p[parent1I].Genotype, p[parent2I].Genotype)
		pt := GrowPhenotype(gt, info)
		p[g] = &Agent{Genotype: gt, Phenotype: pt}
	}
}

func NewPopulation(gts []Genotype, info ActivationInfo) Population {
	p := make(Population, len(gts))
	for i := range p {
		pt := GrowPhenotype(gts[i], info)
		p[i] = &Agent{Genotype: gts[i], Phenotype: pt}
	}
	return p
}
