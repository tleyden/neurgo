package neurgo

import (
	"math"
	"math/rand"
)

type StochasticHillClimber struct {
	currentCandidate *NeuralNetwork
	currentOptimal   *NeuralNetwork
}

func (shc *StochasticHillClimber) Train(neuralNet *NeuralNetwork, examples []*TrainingSample) *NeuralNetwork {

	fittestNeuralNet := neuralNet

	// Apply NN to problem and save fitness
	fitness := fittestNeuralNet.Fitness(examples)

	if fitness > FITNESS_THRESHOLD {
		return fittestNeuralNet
	}

	for {

		// Save the genotype
		candidateNeuralNet := fittestNeuralNet.Copy()

		// Perturb synaptic weights and biases
		shc.perturbParameters(candidateNeuralNet)

		// Re-Apply NN to problem
		candidateFitness := candidateNeuralNet.Fitness(examples)

		// If the fitness of the perturbed NN is higher, discard original NN and keep
		// the new.  If the fitness of original is higher, discard perturbed and keep
		// the old.
		if candidateFitness > fitness {
			fittestNeuralNet = candidateNeuralNet
			fitness = candidateFitness
		}

		if candidateFitness > FITNESS_THRESHOLD {
			break
		}

	}

	return fittestNeuralNet

}

// 1. Each neuron in the neural net (weight or bias) will be chosen for perturbation
//    with a probability of 1/sqrt(nn_size)
// 2. Within the chosen neuron, the weights which will be perturbed will be chosen
//    with probability of 1/sqrt(parameters_size)
// 3. The intensity of the parameter perturbation will chosen with uniform distribution
//    of -pi and pi
func (shc *StochasticHillClimber) perturbParameters(neuralNet *NeuralNetwork) {

	// pick the neurons to perturb (at least one)
	neurons := shc.chooseNeuronsToPerturb(neuralNet)

	for _, neuron := range neurons {
		shc.perturbNeuron(neuron)
	}

}

func (shc *StochasticHillClimber) chooseNeuronsToPerturb(neuralNet *NeuralNetwork) []*Node {

	neuronsToPerturb := make([]*Node, 0)

	// choose some random neurons to perturb.  we need at least one, so
	// keep looping until we've chosen at least one
	didChooseNeuron := false
	for {

		probability := shc.nodePerturbProbability(neuralNet)
		neurons := neuralNet.neurons()
		for _, neuronNode := range neurons {
			if rand.Float64() < probability {
				neuronsToPerturb = append(neuronsToPerturb, neuronNode)
				didChooseNeuron = true
			}
		}

		if didChooseNeuron {
			break
		}

	}
	return neuronsToPerturb

}

func (shc *StochasticHillClimber) nodePerturbProbability(neuralNet *NeuralNetwork) float64 {
	neurons := neuralNet.neurons()
	numNeurons := len(neurons)
	return 1 / math.Sqrt(float64(numNeurons))
}

func (shc *StochasticHillClimber) perturbNeuron(node *Node) {

	probability := shc.parameterPerturbProbability(node)

	// keep trying until we've perturbed at least one parameter
	for {
		didPerturbWeight := false
		for _, cxn := range node.inbound {
			didPerturbWeight = shc.possiblyPerturbConnection(cxn, probability)
		}

		didPerturbBias := shc.possiblyPerturbBias(node, probability)

		// did we perturb anything?  if so, we're done
		if didPerturbWeight || didPerturbBias {
			break
		}

	}

}

func (shc *StochasticHillClimber) parameterPerturbProbability(node *Node) float64 {
	numWeights := 0
	for _, connection := range node.inbound {
		numWeights += len(connection.weights)
	}
	return 1 / math.Sqrt(float64(numWeights))
}

func (shc *StochasticHillClimber) possiblyPerturbConnection(cxn *connection, probability float64) bool {

	didPerturb := false
	for j, weight := range cxn.weights {
		if rand.Float64() < probability {
			perturbedWeight := shc.perturbParameter(weight)
			cxn.weights[j] = perturbedWeight
			didPerturb = true
		}
	}
	return didPerturb

}

func (shc *StochasticHillClimber) possiblyPerturbBias(node *Node, probability float64) bool {
	didPerturb := false
	hasBias := node.processor.hasBias()
	if hasBias && rand.Float64() < probability {
		bias := node.processor.bias()
		perturbedBias := shc.perturbParameter(bias)
		node.processor.setBias(perturbedBias)
		didPerturb = true
	}
	return didPerturb
}

func (shc *StochasticHillClimber) perturbParameter(parameter float64) float64 {

	parameter += RandomInRange(-3.14, 3.14)
	return parameter

}
