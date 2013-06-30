package neurgo

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/go.assert"
	"testing"
)

func simpleNetwork() *NeuralNetwork {

	// create network nodes
	neuronProcessor1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuronProcessor2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuron1 := &Node{Name: "neuron1", processor: neuronProcessor1}
	neuron2 := &Node{Name: "neuron2", processor: neuronProcessor2}
	sensor := &Node{Name: "sensor", processor: &Sensor{}}
	actuator := &Node{Name: "actuator", processor: &Actuator{}}

	// connect nodes together
	weights := []float64{20, 20, 20, 20, 20}
	sensor.ConnectBidirectionalWeighted(neuron1, weights)
	sensor.ConnectBidirectionalWeighted(neuron2, weights)
	neuron1.ConnectBidirectional(actuator)
	neuron2.ConnectBidirectional(actuator)

	// create neural network
	sensors := []*Node{sensor}
	actuators := []*Node{actuator}
	neuralNet := &NeuralNetwork{sensors: sensors, actuators: actuators}

	return neuralNet
}

func TestNetworkVerify(t *testing.T) {

	neuralNet := simpleNetwork()

	// inputs + expected outputs
	examples := []*TrainingSample{{SampleInputs: [][]float64{[]float64{1, 1, 1, 1, 1}}, ExpectedOutputs: [][]float64{[]float64{110, 110}}}}

	// verify neural network
	verified := neuralNet.Verify(examples)
	assert.True(t, verified)

	// make sure injectors/wiretaps have been removed
	assert.Equals(t, len(neuralNet.sensors[0].inbound), 0)
	assert.Equals(t, len(neuralNet.actuators[0].outbound), 0)

}

func TestNetworkFitness(t *testing.T) {

	neuralNet := simpleNetwork()

	// inputs + expected outputs
	examples := []*TrainingSample{{SampleInputs: [][]float64{[]float64{1, 1, 1, 1, 1}}, ExpectedOutputs: [][]float64{[]float64{110, 110}}}}

	// get network fitness
	fitness := neuralNet.Fitness(examples)
	assert.True(t, fitness > 10000000)

	// inputs + crazy outputs
	badExamples := []*TrainingSample{{SampleInputs: [][]float64{[]float64{1, 1, 1, 1, 1}}, ExpectedOutputs: [][]float64{[]float64{-1000, -1000}}}}

	lowFitness := neuralNet.Fitness(badExamples)
	assert.True(t, equalsWithMaxDelta(lowFitness, 0.0, .01))

}

func TestXnorTwoSensorNetwork(t *testing.T) {

	// create network nodes
	n1_processor := &Neuron{Bias: 0, ActivationFunction: identity_activation}
	input_neuron1 := &Node{Name: "input_neuron1", processor: n1_processor}

	n2_processor := &Neuron{Bias: 0, ActivationFunction: identity_activation}
	input_neuron2 := &Node{Name: "input_neuron2", processor: n2_processor}

	hn1_processor := &Neuron{Bias: -30, ActivationFunction: Sigmoid}
	hidden_neuron1 := &Node{Name: "hidden_neuron1", processor: hn1_processor}

	hn2_processor := &Neuron{Bias: 10, ActivationFunction: Sigmoid}
	hidden_neuron2 := &Node{Name: "hidden_neuron2", processor: hn2_processor}

	outn_processor := &Neuron{Bias: -10, ActivationFunction: Sigmoid}
	output_neuron := &Node{Name: "output_neuron", processor: outn_processor}

	sensor1 := &Node{Name: "sensor1", processor: &Sensor{}}
	sensor2 := &Node{Name: "sensor2", processor: &Sensor{}}
	actuator := &Node{Name: "actuator", processor: &Actuator{}}

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
	sensors := []*Node{sensor1, sensor2}
	actuators := []*Node{actuator}
	neuralNet := &NeuralNetwork{sensors: sensors, actuators: actuators}

	// inputs + expected outputs
	examples := []*TrainingSample{

		// TODO: how to wrap this?
		{SampleInputs: [][]float64{[]float64{0}, []float64{1}}, ExpectedOutputs: [][]float64{[]float64{0}}},
		{SampleInputs: [][]float64{[]float64{1}, []float64{1}}, ExpectedOutputs: [][]float64{[]float64{1}}},
		{SampleInputs: [][]float64{[]float64{1}, []float64{0}}, ExpectedOutputs: [][]float64{[]float64{0}}},
		{SampleInputs: [][]float64{[]float64{0}, []float64{0}}, ExpectedOutputs: [][]float64{[]float64{1}}}}

	// verify neural network
	verified := neuralNet.Verify(examples)
	assert.True(t, verified)

}

func xnorTrainingSamples() []*TrainingSample {

	// inputs + expected outputs
	examples := []*TrainingSample{

		// TODO: how to wrap this?
		{SampleInputs: [][]float64{[]float64{0, 1}}, ExpectedOutputs: [][]float64{[]float64{0}}},
		{SampleInputs: [][]float64{[]float64{1, 1}}, ExpectedOutputs: [][]float64{[]float64{1}}},
		{SampleInputs: [][]float64{[]float64{1, 0}}, ExpectedOutputs: [][]float64{[]float64{0}}},
		{SampleInputs: [][]float64{[]float64{0, 0}}, ExpectedOutputs: [][]float64{[]float64{1}}}}

	return examples

}

func TestXnorCondensedNetwork(t *testing.T) {

	// identical to TestXnorNetwork, but uses single sensor with vector outputs, removes
	// the input layer neurons which are useless

	neuralNet := XnorCondensedNetwork()

	// inputs + expected outputs
	examples := xnorTrainingSamples()

	// verify neural network
	verified := neuralNet.Verify(examples)
	assert.True(t, verified)

}

func TestUniqueNodeMap(t *testing.T) {
	neuralNet := XnorCondensedNetwork()
	nodes := neuralNet.uniqueNodeMap()
	assert.Equals(t, len(nodes), 5)
}

func TestGetNeurons(t *testing.T) {
	neuralNet := XnorCondensedNetwork()
	neurons := neuralNet.Neurons()
	assert.Equals(t, len(neurons), 3)
}

func TestShutdown(t *testing.T) {

	neuralNet := XnorCondensedNetwork()
	examples := xnorTrainingSamples()
	for i := 0; i < 25; i++ {
		verify := neuralNet.Verify(examples)
		assert.True(t, verify)
	}

}

func TestNextLayerNodes(t *testing.T) {
	neuralNet := XnorCondensedNetwork()
	sensors := neuralNet.sensors
	nextLayerNodes := neuralNet.nextLayerNodes(sensors)
	assert.Equals(t, len(nextLayerNodes), 2)
}

func TestNodesByLayer(t *testing.T) {
	neuralNet := XnorCondensedNetwork()
	byLayer := neuralNet.NodesByLayer()
	assert.Equals(t, len(byLayer[0]), 1)
	assert.Equals(t, len(byLayer[1]), 2)
	assert.Equals(t, len(byLayer[2]), 1)
	assert.Equals(t, len(byLayer[3]), 1)
}

func TestNumLayers(t *testing.T) {
	neuralNet := XnorCondensedNetwork()
	assert.Equals(t, neuralNet.NumLayers(), 2)
}

func TestCopy(t *testing.T) {

	neuralNet := XnorCondensedNetwork()
	neuralNetCopy := neuralNet.Copy()

	assert.NotEquals(t, neuralNet, neuralNetCopy)
	assert.Equals(t, len(neuralNet.sensors), len(neuralNetCopy.sensors))
	assert.NotEquals(t, neuralNet.sensors[0], neuralNetCopy.sensors[0])
	assert.Equals(t, neuralNet.sensors[0].Name, neuralNetCopy.sensors[0].Name)
	assert.Equals(t, len(neuralNet.actuators), len(neuralNetCopy.actuators))
	assert.NotEquals(t, neuralNet.actuators[0], neuralNetCopy.actuators[0])

	assert.Equals(t, len(neuralNet.sensors[0].outbound), len(neuralNetCopy.sensors[0].outbound))
	assert.NotEquals(t, neuralNet.sensors[0].outbound[0], neuralNetCopy.sensors[0].outbound[0])

	assert.False(t, neuralNetCopy.sensors[0].outbound[0].channel == nil)
	assert.Equals(t, len(neuralNet.actuators[0].inbound), len(neuralNetCopy.actuators[0].inbound))

	assert.Equals(t, len(neuralNetCopy.sensors[0].outbound[0].other.inboundConnections()), len(neuralNet.sensors[0].outbound[0].other.inboundConnections()))

	assert.True(t, neuralNetCopy.sensors[0].outbound[0].channel == neuralNetCopy.sensors[0].outbound[0].other.inboundConnections()[0].channel)

	assert.NotEquals(t, neuralNet.actuators[0].inbound[0], neuralNetCopy.actuators[0].inbound[0])
	assert.Equals(t, len(neuralNetCopy.sensors[0].outbound[0].other.inboundConnections()[0].weights), len(neuralNet.sensors[0].outbound[0].other.inboundConnections()[0].weights))

	otherNeuron := neuralNet.sensors[0].outbound[0].other.processor.(*Neuron)
	otherNeuronCopy := neuralNetCopy.sensors[0].outbound[0].other.processor.(*Neuron)
	assert.Equals(t, otherNeuron.Bias, otherNeuronCopy.Bias)
	assert.Equals(t, otherNeuron.ActivationFunction(1), otherNeuronCopy.ActivationFunction(1))

	assert.True(t, neuralNetCopy.sensors[0].processor != nil)
	assert.True(t, neuralNetCopy.actuators[0].processor != nil)

	outputNeuron := neuralNetCopy.actuators[0].inbound[0].other
	assert.Equals(t, len(outputNeuron.outbound), 1)

	nnJson, _ := json.Marshal(neuralNet)
	nnJsonString := fmt.Sprintf("%s", nnJson)

	nnCopyJson, _ := json.Marshal(neuralNetCopy)
	nnCopyJsonString := fmt.Sprintf("%s", nnCopyJson)
	assert.Equals(t, nnJsonString, nnCopyJsonString)

	examples := xnorTrainingSamples()
	verified := neuralNetCopy.Verify(examples)
	assert.True(t, verified)

}
