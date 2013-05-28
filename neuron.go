
package neurgo

type activationFunction func(float32) float32

type Neuron struct {
	Bias float32
	ActivationFunction activationFunction
}

// Connectable interface implementations

func (neuron *Neuron) Connect_with_weights(target NeuralNode, weights []float32) {

}

func (neuron *Neuron) Connect(target NeuralNode) {

}

// NeuralNode interface implementations

func (neuron *Neuron) DoSomething() {

}

// Methods

func (neuron *Neuron) Run() {

}
