package neurgo

type TrainingSample struct {

	// for each sensor in the network, provide a sample input vector
	SampleInputs [][]float64

	// for each actuator in the network, provide an expected output vector
	ExpectedOutputs [][]float64
}

type Trainer interface {
	Train(neuralNet *NeuralNetwork, examples []*TrainingSample) *NeuralNetwork
}
