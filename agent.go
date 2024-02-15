package goevo

type Agent struct {
	Genotype *Genotype
	Fitness  float64
}

func NewAgent(gt *Genotype) *Agent {
	return &Agent{Genotype: gt}
}
