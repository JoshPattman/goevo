package goevo

import "math"

type GeneticDistanceFunc func(*Genotype, *Genotype) float64

// Returns a function that computes the genetic distance between the two genotypes.
// This is disjoint * number_of_disjoint_genes + matching * total_matching_synapse_weight_diff.
// The values used in the original paper are disjoint=1, matching=0.4.
// The returned function can be used in Speciate().
func ClassicGeneticDistance(disjoint, matching float64) GeneticDistanceFunc {
	return func(g1, g2 *Genotype) float64 {
		numMatching := 0.0
		totalWDiff := 0.0
		for sid1, s1 := range g1.Synapses {
			if s2, ok := g2.Synapses[sid1]; ok {
				// Matching
				numMatching++
				totalWDiff += math.Abs(s1.Weight - s2.Weight)
			}
		}
		numDisjoint := (float64(len(g1.Synapses)) - numMatching) + (float64(len(g2.Synapses)) - numMatching)

		return disjoint*numDisjoint + matching*totalWDiff
	}
}
