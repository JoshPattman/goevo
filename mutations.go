package goevo

import (
	"math/rand"
)

type GenotypeMutator struct {
	MaxNewSynapseValue      float64
	MaxSynapseMutationValue float64
}

func (gm *GenotypeMutator) GrowRandomSynapse(g *GenotypeSlow, counter InnovationCounter) {
	inps, hids, outs := g.GetNodeTypeCounts()
	for i := 0; i < 10; i++ {
		ra := randRange(0, inps+hids)
		s := ra + 1
		if s < inps {
			s = inps
		}
		rb := randRange(s, inps+hids+outs)
		w := (rand.Float64()*2 - 1) * gm.MaxNewSynapseValue
		if g.CreateConnection(g.Layers[ra].ID, g.Layers[rb].ID, w, counter) {
			return
		}
	}
}
func (gm *GenotypeMutator) GrowRandomRecurrentSynapse(g *GenotypeSlow, counter InnovationCounter) {
	inps, hids, outs := g.GetNodeTypeCounts()
	for i := 0; i < 10; i++ {
		ra := randRange(0, inps+hids)
		s := ra + 1
		if s < inps {
			s = inps
		}
		rb := randRange(s, inps+hids+outs)
		w := (rand.Float64()*2 - 1) * gm.MaxNewSynapseValue
		if g.CreateRecurrentConnection(g.Layers[rb].ID, g.Layers[ra].ID, w, counter) {
			return
		}
	}
}

func (gm *GenotypeMutator) GrowRandomNode(g *GenotypeSlow, counter InnovationCounter) {
	cons := make([]*ConnectionGene, len(g.Connections))
	consI := 0
	for i := range g.Connections {
		if g.Connections[i].Enabled && !g.Connections[i].Recurrent {
			cons[consI] = g.Connections[i]
			consI++
		}
	}
	cons = cons[:consI]
	if len(cons) == 0 {
		return
	}
	si := randRange(0, len(cons))
	g.CreateNode(cons[si].ID, counter)
}

func (gm *GenotypeMutator) MutateRandomConnection(g *GenotypeSlow) {
	cons := make([]*ConnectionGene, len(g.Connections))
	consI := 0
	for i := range g.Connections {
		if g.Connections[i].Enabled {
			cons[consI] = g.Connections[i]
			consI++
		}
	}
	cons = cons[:consI]
	if len(cons) == 0 {
		return
	}
	si := randRange(0, len(cons))
	cons[si].Weight += (rand.Float64()*2 - 1) * gm.MaxSynapseMutationValue
}
