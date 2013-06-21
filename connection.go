package neurgo

import (
	"encoding/json"
)

type connection struct {
	other   *Node
	channel VectorChannel
	closing chan bool
	weights []float64
}

func (cxn *connection) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Weights []float64 `json:"weights"`
			Other   string    `json:"other"`
		}{
			Weights: cxn.weights,
			Other:   cxn.other.String(),
		})

}
