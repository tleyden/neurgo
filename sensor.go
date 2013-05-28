
package neurgo

type Sensor struct {
	InputChannel chan []float32
}

// Connectable interface implementations

func (sensor *Sensor) Connect_with_weights(target Connectable, weights []float32) {

	// 

}

func (sensor *Sensor) Connect(target Connectable) {

}

// NeuralNode interface implementations

func (sensor *Sensor) DoSomething() {

}

// Methods

func (sensor *Sensor) Run() {

}

func (sensor *Sensor) SendInput(input []float32) {
	sensor.InputChannel <- input
}
