package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
)


func TestNetworkVerify(t *testing.T) {

	// create network nodes
	neuron1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}  
	neuron2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	sensor := &Sensor{}
	actuator := &Actuator{}

	// give nodes names
	sensor.Name = "sensor"
	actuator.Name = "actuator"
	neuron1.Name = "neuron1"
	neuron2.Name = "neuron2"

	// connect nodes together 
	weights := []float64{20,20,20,20,20}
	sensor.ConnectBidirectionalWeighted(neuron1, weights)
	sensor.ConnectBidirectionalWeighted(neuron2, weights)
	neuron1.ConnectBidirectional(actuator)
	neuron2.ConnectBidirectional(actuator)

	// inputs + expected outputs
	examples := []*TrainingSample{{sampleInputs: [][]float64{[]float64{1,1,1,1,1}}, expectedOutputs: [][]float64{[]float64{110,110}}}}

	// create neural network
	sensors := []*Sensor{sensor}	
	actuators := []*Actuator{actuator}
	neuralNet := &NeuralNetwork{sensors: sensors, actuators: actuators}

	// spinup node goroutines
	signallers := []Signaller{neuron1, neuron2, sensor, actuator}
	for _, signaller := range signallers {
		go Run(signaller)
	}

	// verify neural network
	verified := neuralNet.Verify(examples)
	assert.True(t, verified)

	// make sure injectors/wiretaps have been removed
	assert.Equals(t, len(sensor.inbound), 0)
	assert.Equals(t, len(actuator.outbound), 0)
	

}

func TestXnorNetwork(t *testing.T) {

	// create network nodes
	input_neuron1 := &Neuron{Bias: 0, ActivationFunction: identity_activation}   
	input_neuron2 := &Neuron{Bias: 0, ActivationFunction: identity_activation}  
	hidden_neuron1 := &Neuron{Bias: -30, ActivationFunction: sigmoid}  
	hidden_neuron2 := &Neuron{Bias: 10, ActivationFunction: sigmoid}  
	output_neuron := &Neuron{Bias: -10, ActivationFunction: sigmoid}  
	sensor1 := &Sensor{}
	sensor2 := &Sensor{}
	actuator := &Actuator{}

	// give names to network nodes
	sensor1.Name = "sensor1"
	sensor2.Name = "sensor2"
	input_neuron1.Name = "input_neuron1"
	input_neuron2.Name = "input_neuron2"
	hidden_neuron1.Name = "hidden_neuron1"
	hidden_neuron2.Name = "hidden_neuron2"
	output_neuron.Name = "output_neuron"
	actuator.Name = "actuator"

	// connect nodes together 
	sensor1.ConnectBidirectionalWeighted(input_neuron1, []float64{1})
	sensor2.ConnectBidirectionalWeighted(input_neuron2, []float64{1})
	input_neuron1.ConnectBidirectionalWeighted(hidden_neuron1, []float64{20})
	input_neuron2.ConnectBidirectionalWeighted(hidden_neuron1, []float64{20})
	input_neuron1.ConnectBidirectionalWeighted(hidden_neuron2, []float64{-20})
	input_neuron2.ConnectBidirectionalWeighted(hidden_neuron2, []float64{-20})
	hidden_neuron1.ConnectBidirectionalWeighted(output_neuron, []float64{20})
	hidden_neuron2.ConnectBidirectionalWeighted(output_neuron, []float64{20})
	output_neuron.ConnectBidirectional(actuator)

	// create neural network
	sensors := []*Sensor{sensor1, sensor2}	
	actuators := []*Actuator{actuator}
	neuralNet := &NeuralNetwork{sensors: sensors, actuators: actuators}

	// inputs + expected outputs
	examples := []*TrainingSample{

		// TODO: how to wrap this?
		{sampleInputs: [][]float64{[]float64{0},[]float64{1}}, expectedOutputs: [][]float64{[]float64{0}}},
		{sampleInputs: [][]float64{[]float64{1},[]float64{1}}, expectedOutputs: [][]float64{[]float64{1}}},
		{sampleInputs: [][]float64{[]float64{1},[]float64{0}}, expectedOutputs: [][]float64{[]float64{0}}},
		{sampleInputs: [][]float64{[]float64{0},[]float64{0}}, expectedOutputs: [][]float64{[]float64{1}}}}


	// spinup node goroutines
	signallers := []Signaller{input_neuron1, input_neuron2, hidden_neuron1, hidden_neuron2, output_neuron, sensor1, sensor2, actuator}
	for _, signaller := range signallers {
		go Run(signaller)
	}

	// verify neural network
	verified := neuralNet.Verify(examples)
	assert.True(t, verified)


}

func TestXnorCondensedNetwork(t *testing.T) {

	// identical to TestXnorNetwork, but uses single sensor with vector outputs, removes 
	// the input layer neurons which are useless

	// create network nodes
	hidden_neuron1 := &Neuron{Bias: -30, ActivationFunction: sigmoid}  
	hidden_neuron2 := &Neuron{Bias: 10, ActivationFunction: sigmoid}  
	output_neuron := &Neuron{Bias: -10, ActivationFunction: sigmoid}  
	sensor := &Sensor{}
	actuator := &Actuator{}

	// give names to network nodes
	sensor.Name = "sensor"
	hidden_neuron1.Name = "hidden_neuron1"
	hidden_neuron2.Name = "hidden_neuron2"
	output_neuron.Name = "output_neuron"
	actuator.Name = "actuator"

	// connect nodes together 
	sensor.ConnectBidirectionalWeighted(hidden_neuron1, []float64{20,20})
	sensor.ConnectBidirectionalWeighted(hidden_neuron2, []float64{-20, -20})
	hidden_neuron1.ConnectBidirectionalWeighted(output_neuron, []float64{20})
	hidden_neuron2.ConnectBidirectionalWeighted(output_neuron, []float64{20})
	output_neuron.ConnectBidirectional(actuator)

	// create neural network
	sensors := []*Sensor{sensor}	
	actuators := []*Actuator{actuator}
	neuralNet := &NeuralNetwork{sensors: sensors, actuators: actuators}

	// inputs + expected outputs
	examples := []*TrainingSample{

		// TODO: how to wrap this?
		{sampleInputs: [][]float64{[]float64{0, 1}}, expectedOutputs: [][]float64{[]float64{0}}},
		{sampleInputs: [][]float64{[]float64{1, 1}}, expectedOutputs: [][]float64{[]float64{1}}},
		{sampleInputs: [][]float64{[]float64{1, 0}}, expectedOutputs: [][]float64{[]float64{0}}},
		{sampleInputs: [][]float64{[]float64{0, 0}}, expectedOutputs: [][]float64{[]float64{1}}}}

	// spinup node goroutines
	signallers := []Signaller{sensor, hidden_neuron1, hidden_neuron2, output_neuron, actuator}
	for _, signaller := range signallers {
		go Run(signaller)
	}

	// verify neural network
	verified := neuralNet.Verify(examples)
	assert.True(t, verified)


}
