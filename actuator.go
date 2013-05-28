
package neurgo

type Actuator struct {
	OutputChannel chan float32
}

// Connectable interface implementations

func (actuator *Actuator) Connect_with_weights(target NeuralNode, weights []float32) {
	panic("Actuator does not support outbound connections")
}

func (actuator *Actuator) Connect(target NeuralNode) {
	panic("Actuator does not support outbound connections")
}

// NeuralNode interface implementations

func (actuator *Actuator) DoSomething() {

}

// Methods

func (actuator *Actuator) Run() {
	actuator.OutputChannel <- 99.0
}

func (actuator *Actuator) ReadOutput() float32 {
	return <- actuator.OutputChannel
}
