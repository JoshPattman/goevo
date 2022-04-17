package goevo

import (
	"math/rand"
)

type GenotypeMutator struct {
	MaxNewSynapseValue      float64
	MaxSynapseMutationValue float64
}

func (gm *GenotypeMutator) GrowRandomSynapse(g Genotype, counter InnovationCounter) {
	inps, hids, outs := g.GetNumNodes()
	for i := 0; i < 10; i++ {
		ra := randRange(0, inps+hids)
		s := ra + 1
		if s < inps {
			s = inps
		}
		rb := randRange(s, inps+hids+outs)
		w := (rand.Float64()*2 - 1) * gm.MaxNewSynapseValue
		nida, _ := g.GetNodeIDAtLayer(ra)
		nidb, _ := g.GetNodeIDAtLayer(rb)
		if _, err := g.ConnectNodes(nida, nidb, w, counter); err == nil {
			return
		}
	}
}
func (gm *GenotypeMutator) GrowRandomRecurrentSynapse(g Genotype, counter InnovationCounter) {
	inps, hids, outs := g.GetNumNodes()
	for i := 0; i < 10; i++ {
		ra := randRange(0, inps+hids)
		s := ra + 1
		if s < inps {
			s = inps
		}
		rb := randRange(s, inps+hids+outs)
		w := (rand.Float64()*2 - 1) * gm.MaxNewSynapseValue
		nida, _ := g.GetNodeIDAtLayer(ra)
		nidb, _ := g.GetNodeIDAtLayer(rb)
		if _, err := g.ConnectNodes(nidb, nida, w, counter); err == nil {
			return
		}
	}
}

func (gm *GenotypeMutator) MutateRandomActivation(g Genotype, acs []string) {
	a := acs[randRange(0, len(acs))]
	in, hid, out := g.GetNumNodes()
	n := randRange(0, in+out+hid)
	nid, _ := g.GetNodeIDAtLayer(n)
	g.SetActivation(nid, a)
}

func (gm *GenotypeMutator) GrowRandomNode(g Genotype, acs []string, counter InnovationCounter) {
	gCons := g.GetAllConnectionIDs()
	cons := make([]ConnectionID, len(gCons))
	consI := 0
	for i := range gCons {
		r, _ := g.IsConnectionRecurrent(gCons[i])
		if !r {
			cons[consI] = gCons[i]
			consI++
		}
	}
	cons = cons[:consI]
	if len(cons) == 0 {
		return
	}
	si := randRange(0, len(cons))
	n, _, _ := g.CreateNodeOn(cons[si], counter)
	g.SetActivation(n, acs[randRange(0, len(acs))])
}

func (gm *GenotypeMutator) MutateRandomConnection(g Genotype) {
	gCons := g.GetAllConnectionIDs()
	cons := make([]ConnectionID, len(gCons))
	consI := 0
	for i := range gCons {
		r, _ := g.IsConnectionRecurrent(gCons[i])
		if !r {
			cons[consI] = gCons[i]
			consI++
		}
	}
	cons = cons[:consI]
	if len(cons) == 0 {
		return
	}
	si := randRange(0, len(cons))
	g.MutateConnectionWeight(cons[si], (rand.Float64()*2-1)*gm.MaxSynapseMutationValue)
}
