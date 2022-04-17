package goevo

import "math/rand"

type GenotypeCrossover struct{}

func (c *GenotypeCrossover) CrossoverSimple(g1, g2 *Genotype) *Genotype {
	g := CopyGenotype(g1)
	for sid, _ := range g.Connections {
		c2 := g2.GetConnection(sid)
		if c2 != nil {
			if rand.Float64() > 0.5 {
				g.Connections[sid].Weight = c2.Weight
			}
		}
	}
	return g
}
