package neurgo

import (
	"encoding/json"
	"log"
)

type Sensor struct {
}

func (sensor *Sensor) hasBias() bool {
	return false
}

func (sensor *Sensor) bias() float64 {
	panic("Sensors don't have bias parameter")
}

func (sensor *Sensor) setBias(newBias float64) {
	panic("Sensors don't have bias parameter")
}

func (sensor *Sensor) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Type string `json:"type"`
		}{
			Type: "Sensor",
		})
}

func (sensor *Sensor) copy() SignalProcessor {
	sensorCopy := &Sensor{}
	return sensorCopy
}

func (sensor *Sensor) waitCanPropagate(node *Node) (isShutdown bool) {
	if len(node.inbound) > 1 { // FIXME: data race #1
		log.Panicf("%v has more than one inbound, this is unexpected", node)
	}

	if len(node.inbound) == 0 {
		isShutdown = node.waitForInboundChannel()
	} else {
		isShutdown = false
	}
	return
}

func (sensor *Sensor) propagateSignal(node *Node) bool {

	// this implemenation is just a stub which makes it easy to test for now.
	// at some point, sensors will act as proxies to real virtual sensors,
	// and probably be reading their inputs from sockets.

	var value []float64
	var ok bool

	connection := node.inbound[0]

	select {
	case value = <-connection.channel:
		ok = true
	case <-connection.closing: // skip this connection since its closed
	case <-node.closing:
		return true
	}

	if ok {
		node.scatterOutput(value)
	}

	return false

}
