package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
)


// create netwwork with topology capable of solving XNOR, but which 
// has not been trained yet
func xnorNetworkUntrained() *NeuralNetwork {

	// create network nodes
	hn1_processor := &Neuron{Bias: 0, ActivationFunction: sigmoid}  
	hidden_neuron1 := &Node{Name: "hidden_neuron1", processor: hn1_processor}
	
	hn2_processor := &Neuron{Bias: 0, ActivationFunction: sigmoid}  
	hidden_neuron2 := &Node{Name: "hidden_neuron2", processor: hn2_processor}

	outn_processor := &Neuron{Bias: 0, ActivationFunction: sigmoid}  
	output_neuron := &Node{Name: "output_neuron", processor: outn_processor}

	sensor := &Node{Name: "sensor", processor: &Sensor{}}
	actuator := &Node{Name: "actuator", processor: &Actuator{}}

	// connect nodes together 
	sensor.ConnectBidirectionalWeighted(hidden_neuron1, []float64{0, 0})
	sensor.ConnectBidirectionalWeighted(hidden_neuron2, []float64{0, 0})
	hidden_neuron1.ConnectBidirectionalWeighted(output_neuron, []float64{0})
	hidden_neuron2.ConnectBidirectionalWeighted(output_neuron, []float64{0})
	output_neuron.ConnectBidirectional(actuator)

	// create neural network
	sensors := []*Node{sensor}	
	actuators := []*Node{actuator}
	neuralNet := &NeuralNetwork{sensors: sensors, actuators: actuators}

	// spinup node goroutines
	nodes := []*Node{sensor, hidden_neuron1, hidden_neuron2, output_neuron, actuator}
	for _, node := range nodes {
		go Run(node.processor, node)
	}
	
	return neuralNet

}

func TestWeightTraining(t *testing.T) {

	// training set
	examples := []*TrainingSample{
		// TODO: how to wrap this?
		{sampleInputs: [][]float64{[]float64{0, 1}}, expectedOutputs: [][]float64{[]float64{0}}},
		{sampleInputs: [][]float64{[]float64{1, 1}}, expectedOutputs: [][]float64{[]float64{1}}},
		{sampleInputs: [][]float64{[]float64{1, 0}}, expectedOutputs: [][]float64{[]float64{0}}},
		{sampleInputs: [][]float64{[]float64{0, 0}}, expectedOutputs: [][]float64{[]float64{1}}}}

	// create netwwork with topology capable of solving XNOR
	neuralNet := xnorNetworkUntrained()

	// verify it can not yet solve the training set (since training would be useless in that case)
	verified := neuralNet.Verify(examples)
	assert.False(t, verified)

	// TODO - shutdown the network so we can re-use it

	// create stochastic hill climber trainer

	// train the network 


	// verify it can now solve the training set
	verified = neuralNet.Verify(examples)
	assert.True(t, verified)


}
