package neurgo

import (
	"testing"
	"fmt"
)

func Test(t *testing.T) {

	fmt.Println("test is running!")

	neuron := &Neuron{}
	sensor := &Sensor{}

	// connect nodes together
	sensor.Connect_with_weights(neuron, []float32{20,20,20,20,20})

	// 

	t.Errorf("fake failure")

}
