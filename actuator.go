
package neurgo

type Actuator struct {

}

func (actuator *Actuator) Connect_with_weights(target NeuralNode, weights []float32) {
	panic("Actuator does not support outbound connections")
}

func (actuator *Actuator) Connect(target NeuralNode) {
	panic("Actuator does not support outbound connections")
}

func (actuator *Actuator) DoSomething() {

}

