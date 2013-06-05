
package neurgo

import (
	"fmt"
)

type Actuator struct {
	Node
}

func (actuator *Actuator) propagateSignal() {

	if len(actuator.inbound) > 0 {  

		outputVectorDimension := len(actuator.inbound)
		outputVector := make([]float64,outputVectorDimension) 

		for i, inboundConnection := range actuator.inbound {
			inputVector := <- inboundConnection.channel

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
