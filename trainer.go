package neurgo

type TrainingSample struct {

	// for each sensor in the network, provide a sample input vector
	sampleInputs [][]float64

	// for each actuator in the network, provide an expected output vector
	expectedOutputs [][]float64
}

type Trainer interface {
	Train(neuralNet *NeuralNetwork, examples []*TrainingSample) *NeuralNetwork
}
