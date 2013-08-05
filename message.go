package neurgo

import (
	"fmt"
)

type DataMessage struct {
	SenderId *NodeId
	Inputs   []float64
}

func (dataMessage *DataMessage) String() string {
	return fmt.Sprintf("sender: %v, values: %v",
		dataMessage.SenderId,
		dataMessage.Inputs,
	)
}
