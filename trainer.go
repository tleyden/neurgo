package neurgo

type TrainingSample struct {
	sampleInput []float64
	expectedOutput []float64
}

type StochasticHillClimber struct {
	currentCandidate []*Node
	currentOptimal []*Node
}

type Trainer interface {
	train(examples []*TrainingSample) []*Node 
}

func (shc *StochasticHillClimber) train(examples []*TrainingSample) []*Node {
	return make([]*Node, 1)
}
