package neurgo

import ()

type Sensor struct {
}

func (sensor *Sensor) copy() SignalProcessor {
	sensorCopy := &Sensor{}
	return sensorCopy
}

func (sensor *Sensor) canPropagateSignal(node *Node) bool {
	return len(node.inbound) == 1
}

func (sensor *Sensor) propagateSignal(node *Node) {

	// this implemenation is just a stub which makes it easy to test for now.
	// at some point, sensors will act as proxies to real virtual sensors,
	// and probably be reading their inputs from sockets.

	if value, ok := <-node.inbound[0].channel; ok {
		node.scatterOutput(value)
	}

}
