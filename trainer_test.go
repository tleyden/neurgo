package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"math"
	"testing"
)

// create netwwork with topology capable of solving XNOR, but which
// has not been trained yet
func xnorNetworkUntrained() *NeuralNetwork {

	randBias_1 := RandomInRange(-1*math.Pi, math.Pi)
	randBias_2 := RandomInRange(-1*math.Pi, math.Pi)
	randBias_3 := RandomInRange(-1*math.Pi, math.Pi)

	randWeights_1 := []float64{
		RandomInRange(-1*math.Pi, math.Pi),
		RandomInRange(-1*math.Pi, math.Pi),
	}

	randWeights_2 := []float64{
		RandomInRange(-1*math.Pi, math.Pi),
		RandomInRange(-1*math.Pi, math.Pi),
	}

	randWeights_3 := []float64{
		RandomInRange(-1*math.Pi, math.Pi),
	}

	randWeights_4 := []float64{
		RandomInRange(-1*math.Pi, math.Pi),
	}

	// create network nodes
	hn1_processor := &Neuron{Bias: randBias_1, ActivationFunction: sigmoid}
	hidden_neuron1 := &Node{Name: "hidden_neuron1", processor: hn1_processor}

	hn2_processor := &Neuron{Bias: randBias_2, ActivationFunction: sigmoid}
	hidden_neuron2 := &Node{Name: "hidden_neuron2", processor: hn2_processor}

	outn_processor := &Neuron{Bias: randBias_3, ActivationFunction: sigmoid}
	output_neuron := &Node{Name: "output_neuron", processor: outn_processor}

	sensor := &Node{Name: "sensor", processor: &Sensor{}}
	actuator := &Node{Name: "actuator", processor: &Actuator{}}

	// connect nodes together
	sensor.ConnectBidirectionalWeighted(hidden_neuron1, randWeights_1)
	sensor.ConnectBidirectionalWeighted(hidden_neuron2, randWeights_2)
	hidden_neuron1.ConnectBidirectionalWeighted(output_neuron, randWeights_3)
	hidden_neuron2.ConnectBidirectionalWeighted(output_neuron, randWeights_4)
	output_neuron.ConnectBidirectional(actuator)

	// create neural network
	sensors := []*Node{sensor}
	actuators := []*Node{actuator}
	neuralNet := &NeuralNetwork{sensors: sensors, actuators: actuators}

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

	shc := new(StochasticHillClimber)
	neuralNetTrained := shc.Train(neuralNet, examples)
	// verify it can now solve the training set
	verified = neuralNetTrained.Verify(examples)
	assert.True(t, verified)

}
