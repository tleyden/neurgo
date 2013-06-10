package neurgo

type TrainingSample struct {

	// for each sensor in the network, provide a sample input vector
	sampleInputs [][]float64

	// for each actuator in the network, provide an expected output vector
	expectedOutputs [][]float64  

}

type StochasticHillClimber struct {
	currentCandidate *NeuralNetwork
	currentOptimal *NeuralNetwork
}

type Trainer interface {
	train(neuralNet *NeuralNetwork, examples []*TrainingSample) *NeuralNetwork 
}

func (shc *StochasticHillClimber) train(examples []*TrainingSample) *NeuralNetwork {
	return &NeuralNetwork{}
}
