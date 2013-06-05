
package neurgo

import (
	"fmt"
	"log"
)

type Actuator struct {
	Node
}

func (actuator *Actuator) propagateSignal() {

	log.Printf("%s: Run().  Type: %T", actuator.Name, actuator) // TODO: how do I print the type of this struct w/o using Name field?

	if len(actuator.inbound) > 0 {  

		outputVectorDimension := len(actuator.inbound)
		outputVector := make([]float32,outputVectorDimension) 

		for i, inboundConnection := range actuator.inbound {
			log.Printf("%v reading from channel: %v", actuator.Name, inboundConnection.channel)
			inputVector := <- inboundConnection.channel

			log.Printf("%v got data from channel: %v", actuator.Name, inboundConnection.channel)

			// assert that the neruron feeding this actuator is emitting a vector containing a single value
			if len(inputVector) != 1 {
				message := fmt.Sprintf("%v got invalid input vector: %v from %v", actuator.Name, inputVector, inboundConnection)
				panic(message)
			}

			inputValue := inputVector[0]

			outputVector[i] = inputValue 
		}

		if len(actuator.outbound) > 0 {
			actuator.outbound[0].channel <- outputVector   
		}

	}

}
