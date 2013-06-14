
package neurgo

import (
)

type Sensor struct {
	Node
}

func (sensor *Sensor) canPropagateSignal() bool {
	return len(sensor.inbound) == 1 
}

func (sensor *Sensor) propagateSignal() {

	// this implemenation is just a stub which makes it easy to test for now.
	// at some point, sensors will act as proxies to real virtual sensors,
	// and probably be reading their inputs from sockets.

	if value, ok := <- sensor.inbound[0].channel; ok {
		sensor.scatterOutput(value)
	} 



}



