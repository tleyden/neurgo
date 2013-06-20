package neurgo

type TrainingSample struct {

	// for each sensor in the network, provide a sample input vector
	sampleInputs [][]float64

	// for each actuator in the network, provide an expected output vector
	expectedOutputs [][]float64
}

type StochasticHillClimber struct {
	currentCandidate *NeuralNetwork
	currentOptimal   *NeuralNetwork
}

type Trainer interface {
	train(neuralNet *NeuralNetwork, examples []*TrainingSample) *NeuralNetwork
}

func (shc *StochasticHillClimber) train(neuralNet *NeuralNetwork, examples []*TrainingSample) *NeuralNetwork {

	/*
		// Repeat

		// Apply NN to problem and save fitness
		fitness = neuralNet.FitnessScore(examples)

		// Save the genotype
		currentWinner := neuralNet.Copy()

		// Perturb synaptic weights and biases
		shc.PerturbParameters(neuralNet)

		// Re-Apply NN to problem
		newFitness := neuralNet.FitnessScore(examples)

		// If the fitness of the perturbed NN is higher, discard original NN and keep
		// the new.  If the fitness of original is higher, discard perturbed and keep
		// the old.
		if newFitness < fitness {
			neuralNet.RestoreFrom(currentWinner)
		}

		// Until - acceptable solution is found, or stopping condition reached

		// Return - genotype with the fittest combination of weights
	*/

	return &NeuralNetwork{}
}
