package neurgo

func XnorCortex() *Cortex {

	// create network nodes

	shouldReInit := false

	sensor := &Sensor{
		NodeId:       NewSensorId("sensor", 0.0),
		VectorLength: 2,
	}
	sensor.Init(shouldReInit)

	hiddenNeuron1 := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("hidden-neuron1", 0.25),
		Bias:               -30,
	}
	hiddenNeuron1.Init(shouldReInit)

	hiddenNeuron2 := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("hidden-neuron2", 0.25),
		Bias:               10,
	}
	hiddenNeuron2.Init(shouldReInit)

	outputNeuron := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("output-neuron", 0.35),
		Bias:               -10,
	}
	outputNeuron.Init(shouldReInit)

	actuator := &Actuator{
		NodeId:       NewActuatorId("actuator", 0.5),
		VectorLength: 1,
	}
	actuator.Init(shouldReInit)

	// wire up connections

	sensor.ConnectOutbound(hiddenNeuron1)
	hiddenNeuron1.ConnectInboundWeighted(sensor, []float64{20, 20})

	sensor.ConnectOutbound(hiddenNeuron2)
	hiddenNeuron2.ConnectInboundWeighted(sensor, []float64{-20, -20})

	hiddenNeuron1.ConnectOutbound(outputNeuron)
	outputNeuron.ConnectInboundWeighted(hiddenNeuron1, []float64{20})

	hiddenNeuron2.ConnectOutbound(outputNeuron)
	outputNeuron.ConnectInboundWeighted(hiddenNeuron2, []float64{20})

	outputNeuron.ConnectOutbound(actuator)
	actuator.ConnectInbound(outputNeuron)

	// create cortex

	nodeId := NewCortexId("cortex")

	cortex := &Cortex{
		NodeId: nodeId,
	}
	cortex.SetSensors([]*Sensor{sensor})
	cortex.SetNeurons([]*Neuron{hiddenNeuron1, hiddenNeuron2, outputNeuron})
	cortex.SetActuators([]*Actuator{actuator})

	return cortex

}

func XnorTrainingSamples() []*TrainingSample {

	// inputs + expected outputs
	examples := []*TrainingSample{

		// TODO: how to wrap this?
		{SampleInputs: [][]float64{[]float64{0, 1}}, ExpectedOutputs: [][]float64{[]float64{0}}},
		{SampleInputs: [][]float64{[]float64{1, 1}}, ExpectedOutputs: [][]float64{[]float64{1}}},
		{SampleInputs: [][]float64{[]float64{1, 0}}, ExpectedOutputs: [][]float64{[]float64{0}}},
		{SampleInputs: [][]float64{[]float64{0, 0}}, ExpectedOutputs: [][]float64{[]float64{1}}}}

	return examples

}

func BasicCortex() *Cortex {

	shouldReInit := false

	// create nodes
	sensor := &Sensor{
		NodeId:       NewSensorId("sensor", 0.0),
		VectorLength: 2,
	}
	sensor.Init(shouldReInit)

	neuron := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("neuron", 0.25),
		Bias:               0,
	}
	neuron.Init(shouldReInit)

	actuator := &Actuator{
		NodeId:       NewActuatorId("actuator", 0.5),
		VectorLength: 1,
	}
	actuator.Init(shouldReInit)

	// wire up connections
	sensor.ConnectOutbound(neuron)
	neuron.ConnectInboundWeighted(sensor, []float64{20, 20})
	neuron.ConnectOutbound(actuator)
	actuator.ConnectInbound(neuron)

	// create cortex
	nodeId := NewCortexId("cortex")
	cortex := &Cortex{
		NodeId: nodeId,
	}
	cortex.SetSensors([]*Sensor{sensor})
	cortex.SetNeurons([]*Neuron{neuron})
	cortex.SetActuators([]*Actuator{actuator})

	return cortex

}
