package neurgo

func XnorCortex() *Cortex {

	// create network nodes

	sensor := &Sensor{
		NodeId:       NewSensorId("sensor", 0.0),
		VectorLength: 2,
	}
	sensor.Init()

	hiddenNeuron1 := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("hidden-neuron1", 0.25),
		Bias:               -30,
	}
	hiddenNeuron1.Init()

	hiddenNeuron2 := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("hidden-neuron2", 0.25),
		Bias:               10,
	}
	hiddenNeuron2.Init()

	outputNeuron := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("output-neuron", 0.35),
		Bias:               -10,
	}
	outputNeuron.Init()

	actuator := &Actuator{
		NodeId:       NewActuatorId("actuator", 0.5),
		VectorLength: 1,
	}
	actuator.Init()

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

func XnorCortexUntrained() *Cortex {

	sensorNodeId := NewSensorId("sensor", 0.0)
	hiddenNeuron1NodeId := NewNeuronId("hidden-neuron1", 0.25)
	hiddenNeuron2NodeId := NewNeuronId("hidden-neuron2", 0.25)
	outputNeuronNodeIde := NewNeuronId("output-neuron", 0.35)

	actuatorNodeId := NewActuatorId("actuator", 0.5)

	hiddenNeuron1 := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             hiddenNeuron1NodeId,
		Bias:               RandomBias(),
	}
	hiddenNeuron1.Init()

	hiddenNeuron2 := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             hiddenNeuron2NodeId,
		Bias:               RandomBias(),
	}
	hiddenNeuron2.Init()

	outputNeuron := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             outputNeuronNodeIde,
		Bias:               RandomBias(),
	}
	outputNeuron.Init()

	sensor := &Sensor{
		NodeId:       sensorNodeId,
		VectorLength: 2,
	}
	sensor.Init()

	actuator := &Actuator{
		NodeId:       actuatorNodeId,
		VectorLength: 1,
	}
	actuator.Init()

	sensor.ConnectOutbound(hiddenNeuron1)
	hiddenNeuron1.ConnectInboundWeighted(sensor, RandomWeights(2))

	sensor.ConnectOutbound(hiddenNeuron2)
	hiddenNeuron2.ConnectInboundWeighted(sensor, RandomWeights(2))

	hiddenNeuron1.ConnectOutbound(outputNeuron)
	outputNeuron.ConnectInboundWeighted(hiddenNeuron1, RandomWeights(1))

	hiddenNeuron2.ConnectOutbound(outputNeuron)
	outputNeuron.ConnectInboundWeighted(hiddenNeuron2, RandomWeights(1))

	outputNeuron.ConnectOutbound(actuator)
	actuator.ConnectInbound(outputNeuron)

	nodeId := NewCortexId("cortex")

	cortex := &Cortex{
		NodeId:    nodeId,
		Sensors:   []*Sensor{sensor},
		Neurons:   []*Neuron{hiddenNeuron1, hiddenNeuron2, outputNeuron},
		Actuators: []*Actuator{actuator},
	}

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

	// create nodes
	sensor := &Sensor{
		NodeId:       NewSensorId("sensor", 0.0),
		VectorLength: 2,
	}
	sensor.Init()

	neuron := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("neuron", 0.25),
		Bias:               0,
	}
	neuron.Init()

	actuator := &Actuator{
		NodeId:       NewActuatorId("actuator", 0.5),
		VectorLength: 1,
	}
	actuator.Init()

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
