package goevo

import (
	"math/rand"
	"sort"
	"math"
)

type Agent struct {
	GT      *Genotype
	PT      *Phenotype
	Fitness float64
}

type Population []*Agent
type SpeciesSettings struct{
	TargetSpecies int
	SThreshold float64
	SThresholdChange float64
}

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
		if p[parent1I].Fitness < p[parent2I].Fitness {
			c := parent2I
			parent2I = parent1I
			parent1I = c
		}
		gt := f(p[parent1I].GT, p[parent2I].GT)
		pt := GrowPhenotype(gt)
		p[g] = &Agent{GT: gt, PT: pt}
	}
}

type SpeciatedPopulation struct{
	Species []Population
	SpeciesFitnesses []float64
	SpeciesAjFitnesses []float64
	SpeciesOffspring []int
	GlobalFitness float64
	GlobalAjFitness float64
}

func NewSpeciatedPopulation()*SpeciatedPopulation{
	return &SpeciatedPopulation{
		make([]Population, 0),
		make([]float64, 0),
		make([]float64, 0),
		make([]int, 0),
		0,
		0,
	}
}

func (p *Population) RepopulateWithSpecies(s *SpeciesSettings, f func(g1 *Genotype, g2 *Genotype) *Genotype) *SpeciatedPopulation{
	// Speciate population
	species := NewSpeciatedPopulation()
    pCopy := make(Population, len(*p))
    copy(pCopy, *p)
    for len(pCopy) > 0{
        thisSpecies := make(Population, 1)
        var initial *Agent
        initial, pCopy = popI(pCopy, rand.Intn(len(pCopy)))
        thisSpecies[0] = initial
        matched, unmatched := make(Population, 0), make(Population, 0)
        for _, a := range pCopy{
            if a.GT.ApproximateGeneticDistance(initial.GT) < s.SThreshold{
                matched = append(matched, a)
            } else{
                unmatched = append(unmatched, a)
            }
        }
        thisSpecies = append(thisSpecies, matched...)
        pCopy = unmatched
        species.Species = append(species.Species, thisSpecies)
    }
    // Ajusting threshold for future speciations
    if len(species.Species) > s.TargetSpecies{
        s.SThreshold += s.SThresholdChange
    } else if len(species.Species) < s.TargetSpecies{
        s.SThreshold -= s.SThresholdChange
    }
	// Calculating fitnesses etc
	gFitness := 0.0
	gaFitness := 0.0
    for _, p := range species.Species{
        pFitness := 0.0
        for _, m := range p{
            pFitness += m.Fitness
        }
        n := float64(len(p))
		sf := pFitness/n
        species.SpeciesFitnesses = append(species.SpeciesFitnesses, sf)
		species.SpeciesAjFitnesses = append(species.SpeciesAjFitnesses, sf/n)
        gFitness += sf
		gaFitness += sf/n
    }
    gFitness = gFitness / float64(len(*p))
	gaFitness = gaFitness / float64(len(*p))
	species.GlobalFitness = gFitness
	species.GlobalAjFitness = gaFitness
    for pi, p := range species.Species{
        species.SpeciesOffspring = append(species.SpeciesOffspring, int(math.Round((species.SpeciesAjFitnesses[pi] / gFitness) * float64(len(p)))))
    }
	// Performing crossover
	gs := make([]*Genotype, 0)
	for si, s := range species.Species{
		s.Order()
		for c := 0; c < species.SpeciesOffspring[si]; c++{
			// pa is taken from the top quater
			// pb is taken from the top half
			pa := rand.Intn(int(math.Ceil(float64(len(s))/4)))
			pb := rand.Intn(int(math.Ceil(float64(len(s))/2)))
			gs = append(gs, f(s[pa].GT, s[pb].GT))
		}
	}
	pNew := NewPopulation(gs)
	(*p) = pNew
	return species
}

func NewPopulation(gts []*Genotype) Population {
	p := make(Population, len(gts))
	for i := range p {
		pt := GrowPhenotype(gts[i])
		p[i] = &Agent{GT: gts[i], PT: pt}
	}
	return p
}

func popI[T any](ts []T, i int) (T, []T){
    e := ts[i]
    ts[i] = ts[len(ts)-1]
    return e, ts[:len(ts)-1]
}
