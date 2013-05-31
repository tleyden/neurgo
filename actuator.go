
package neurgo

type Actuator struct {
	OutputChannel chan float32
	Node
}

// Connectable interface implementations
/*
func (actuator *Actuator) Connect_with_weights(target Connectable, weights []float32) {
	panic("Actuator does not support outbound connections")
}

func (actuator *Actuator) Connect(target Connectable) {
	panic("Actuator does not support outbound connections")
}
*/

// Methods

func (actuator *Actuator) Run() {
	actuator.OutputChannel <- 99.0
}

func (actuator *Actuator) ReadOutput() float32 {
	return <- actuator.OutputChannel
}
