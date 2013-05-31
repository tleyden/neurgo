
package neurgo

import "fmt"

type Sensor struct {
	InputChannel VectorChannel
	Node
}


// Methods

func (sensor *Sensor) Run() {
	fmt.Println("sensor.Run()")
}

func (sensor *Sensor) SendInput(input []float32) {
	sensor.InputChannel <- input
}
