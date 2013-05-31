
package neurgo

import "fmt"

type Sensor struct {
	InputChannel VectorChannel
	NeuralNode
}


// Methods

func (sensor *Sensor) Run() {
	fmt.Println("sensor.Run()")
}

func (sensor *Sensor) SendInput(input []float32) {
	sensor.InputChannel <- input
}
