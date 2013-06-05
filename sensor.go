
package neurgo

import (
	"fmt"
)

type Sensor struct {
	Node
}

func (sensor *Sensor) propagateSignal() {

	// read from input channel and broadcast to all output channels
	if (len(sensor.inbound) != 1) {
		message := fmt.Sprintf("Sensor (%v) should have exactly one input channel, currently has: %v", sensor.Name, len(sensor.inbound))
		panic(message)
	}
	
	value := <- sensor.inbound[0].channel 

	for _, outboundConnection := range sensor.outbound {
		outboundConnection.channel <- value
	}


}



