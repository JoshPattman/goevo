package goevo

import "math/rand"

type GenotypeCrossover struct{}

func (c *GenotypeCrossover) CrossoverSimple(g1, g2 Genotype) Genotype {
	g := g1.Copy()
	allCons := g.GetAllConnectionIDs()
	for _, cid := range allCons {
		w2, err := g2.GetConnectionWeight(cid)
		if err == nil {
			if rand.Float64() > 0.5 {
				g.SetConnectionWeight(cid, w2)
			}
		}
	}
	return g
}
