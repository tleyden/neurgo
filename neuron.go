
package neurgo

type activationFunction func(float32) float32

type Neuron struct {
	Bias               float32
	ActivationFunction activationFunction
	Node
}
