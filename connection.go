package neurgo

import (
	"encoding/json"
)

type Connection struct {
	other   *Node
	channel VectorChannel
	weights []float64
}

func (cxn *Connection) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Weights []float64 `json:"weights"`
			Other   string    `json:"other"`
		}{
			Weights: cxn.weights,
			Other:   cxn.other.String(),
		})

}

func (cxn *Connection) Weights() []float64 {
	return cxn.weights
}
