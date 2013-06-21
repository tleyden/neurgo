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

	var value []float64
	var ok bool

	connection := node.inbound[0]

	select {
	case value = <-connection.channel:
		ok = true
	case <-connection.closing:
	}

	if ok {
		node.scatterOutput(value)
	}

}
