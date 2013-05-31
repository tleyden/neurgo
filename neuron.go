
package neurgo

import "fmt"

type activationFunction func(float32) float32

type Neuron struct {
	Bias float32
	ActivationFunction activationFunction
	NeuralNode
}


// Methods

func (neuron *Neuron) Run() {
	fmt.Println("neuron.Run()")

	// loop through all the channels

	// read value

	// calculate dot product + bias

	// send to output channels

}
