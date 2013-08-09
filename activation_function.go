package neurgo

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
)

type ActivationFunction func(float64) float64

type EncodableActivation struct {
	Name               string
	ActivationFunction ActivationFunction
}

func (activation *EncodableActivation) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Name string
		}{
			Name: activation.Name,
		})
}

func (activation *EncodableActivation) UnmarshalJSON(bytes []byte) error {

	rawMap := make(map[string]interface{})
	err := json.Unmarshal(bytes, &rawMap)
	if err != nil {
		return err
	}

	// TODO: isn't there an easier / less brittle way to do this??
	var ok bool
	if activation.Name, ok = rawMap["Name"].(string); !ok {
		log.Panicf("Could not unmarshal %v into EncodableActivation", rawMap)
	}

	switch activation.Name {
	case "sigmoid":
		activation.ActivationFunction = Sigmoid
	default:
		log.Panicf("Unknown activation function: %v", activation.Name)
	}

	return nil
}

func (activation *EncodableActivation) String() string {
	return fmt.Sprintf("%v (%v)", activation.Name, activation.ActivationFunction)
}

func Sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Pow(math.E, -1*x))
}

func EncodableSigmoid() *EncodableActivation {
	return &EncodableActivation{
		Name:               "sigmoid",
		ActivationFunction: Sigmoid,
	}
}
