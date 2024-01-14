package goevo

// DynamicDistanceThreshold is a simple helper type to update the distance threshold during a generational loop.
// It simply multiplies or divides the current distance threshold each generation, depending on the number of species.
type DynamicDistanceThreshold struct {
	Threshold     float64
	Multiplier    float64
	TargetSpecies int
}

// NewDynamicDistanceThreshold creates a new DynamicDistanceThreshold.
//
//   - multiplier: What values should the threshold be multiplied by each generation
//   - targetSpecies: How many species should be targeted each generation
func NewDynamicDistanceThreshold(initialDistanceThreshold, multiplier float64, targetSpecies int) *DynamicDistanceThreshold {
	return &DynamicDistanceThreshold{
		Threshold:     initialDistanceThreshold,
		Multiplier:    multiplier,
		TargetSpecies: targetSpecies,
	}
}

// UpdateDistanceThreshold updates the distance threshold for this generation, to make the number of species closer to the target next generation.
//
//   - species: A map of species in this generation, from the Speciate function
func (ddt *DynamicDistanceThreshold) UpdateDistanceThreshold(species map[int][]*Agent) {
	if len(species) > ddt.TargetSpecies {
		ddt.Threshold *= ddt.Multiplier
	} else if len(species) < ddt.TargetSpecies {
		ddt.Threshold /= ddt.Multiplier
	}
}
