package goevo

type DynamicDistanceThreshold struct {
	CurrentDistanceThreshold float64
	Multiplier               float64
	TargetSpecies            int
}

func NewDynamicDistanceThreshold(initialDistanceThreshold, multiplier float64, targetSpecies int) *DynamicDistanceThreshold {
	return &DynamicDistanceThreshold{
		CurrentDistanceThreshold: initialDistanceThreshold,
		Multiplier:               multiplier,
		TargetSpecies:            targetSpecies,
	}
}

func (ddt *DynamicDistanceThreshold) UpdateDistanceThreshold(species map[int][]*Agent) {
	if len(species) > ddt.TargetSpecies {
		ddt.CurrentDistanceThreshold *= ddt.Multiplier
	} else if len(species) < ddt.TargetSpecies {
		ddt.CurrentDistanceThreshold /= ddt.Multiplier
	}
}

func (ddt *DynamicDistanceThreshold) Thresh() float64 {
	return ddt.CurrentDistanceThreshold
}
