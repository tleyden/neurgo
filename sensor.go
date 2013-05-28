
package neurgo

type Sensor struct {
	InputChannel chan []float32
}

// Connectable interface implementations

func (sensor *Sensor) Connect_with_weights(target NeuralNode, weights []float32) {

}

func (sensor *Sensor) Connect(target NeuralNode) {

}

// NeuralNode interface implementations

func (sensor *Sensor) DoSomething() {

}


// Methods

func (sensor *Sensor) Run() {

}
