package goevo

type DynamicDistanceThreshold struct {
	Threshold     float64
	Multiplier    float64
	TargetSpecies int
}

func NewDynamicDistanceThreshold(initialDistanceThreshold, multiplier float64, targetSpecies int) *DynamicDistanceThreshold {
	return &DynamicDistanceThreshold{
		Threshold:     initialDistanceThreshold,
		Multiplier:    multiplier,
		TargetSpecies: targetSpecies,
	}
}

func (ddt *DynamicDistanceThreshold) UpdateDistanceThreshold(species map[int][]*Agent) {
	if len(species) > ddt.TargetSpecies {
		ddt.Threshold *= ddt.Multiplier
	} else if len(species) < ddt.TargetSpecies {
		ddt.Threshold /= ddt.Multiplier
	}
}
