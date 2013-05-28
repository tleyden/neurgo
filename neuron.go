
package neurgo

type activationFunction func(float32) float32

type Neuron struct {
	Bias float32
	ActivationFunction activationFunction
}

func (neuron *Neuron) Connect_with_weights(target NeuralNode, weights []float32) {

}

func (neuron *Neuron) Connect(target NeuralNode) {

}

func (neuron *Neuron) DoSomething() {

}

func (neuron *Neuron) Activate(input float32) float32 {  // just a test method.. ignore this
	return neuron.ActivationFunction(input)
}
