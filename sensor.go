
package neurgo

import (
	"log"
	"fmt"
)

type Sensor struct {
	Node
}

func (sensor *Sensor) propagateSignal() {

	log.Printf("%s: Run()", sensor.Name) 

	// read from input channel and broadcast to all output channels
	if (len(sensor.inbound) != 1) {
		message := fmt.Sprintf("Sensor (%v) should have exactly one input channel, currently has: %v", sensor.Name, len(sensor.inbound))
		panic(message)
	}
	
	value := <- sensor.inbound[0].channel 

	for _, outboundConnection := range sensor.outbound {
		log.Printf("%v sending to channel: %v", sensor.Name, outboundConnection.channel)
		outboundConnection.channel <- value
		log.Printf("%v sent to channel: %v", sensor.Name, outboundConnection.channel)
	}


}



