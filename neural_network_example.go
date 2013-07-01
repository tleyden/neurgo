package neurgo

func XnorTrainingSamples() []*TrainingSample {

	examples := []*TrainingSample{
		// TODO: how to wrap this?
		{SampleInputs: [][]float64{[]float64{0, 1}}, ExpectedOutputs: [][]float64{[]float64{0}}},
		{SampleInputs: [][]float64{[]float64{1, 1}}, ExpectedOutputs: [][]float64{[]float64{1}}},
		{SampleInputs: [][]float64{[]float64{1, 0}}, ExpectedOutputs: [][]float64{[]float64{0}}},
		{SampleInputs: [][]float64{[]float64{0, 0}}, ExpectedOutputs: [][]float64{[]float64{1}}}}

	return examples

}

func XnorCondensedNetwork() *NeuralNetwork {

	// create network nodes
	hn1_processor := &Neuron{Bias: -30, ActivationFunction: Sigmoid}
	hidden_neuron1 := &Node{Name: "hidden_neuron1", processor: hn1_processor}

	hn2_processor := &Neuron{Bias: 10, ActivationFunction: Sigmoid}
	hidden_neuron2 := &Node{Name: "hidden_neuron2", processor: hn2_processor}

	outn_processor := &Neuron{Bias: -10, ActivationFunction: Sigmoid}
	output_neuron := &Node{Name: "output_neuron", processor: outn_processor}

	sensor := &Node{Name: "sensor", processor: &Sensor{}}
	actuator := &Node{Name: "actuator", processor: &Actuator{}}

	// connect nodes together
	sensor.ConnectBidirectionalWeighted(hidden_neuron1, []float64{20, 20})
	sensor.ConnectBidirectionalWeighted(hidden_neuron2, []float64{-20, -20})
	hidden_neuron1.ConnectBidirectionalWeighted(output_neuron, []float64{20})
	hidden_neuron2.ConnectBidirectionalWeighted(output_neuron, []float64{20})
	output_neuron.ConnectBidirectional(actuator)

	// create neural network
	sensors := []*Node{sensor}
	actuators := []*Node{actuator}
	neuralNet := &NeuralNetwork{sensors: sensors, actuators: actuators}

	return neuralNet
}
