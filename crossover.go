package goevo

type GenotypeCrossover struct {
}

/*
func (c *GenotypeCrossover) CrossoverSimple(g1, g2 *Genotype) *Genotype {
	g := CopyGenotype(g1)
	for s := range g.Connections {
		c2, _ := g2.GetConnection(g.Connections[s].ID)
		if c2 != nil {
			if rand.Float64() > 0.5 {
				g.Connections[s] = *c2
			}
		}
	}
	return g
}
*/
