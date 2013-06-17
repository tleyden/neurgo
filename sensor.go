
package neurgo

import (
)

type Sensor struct {
	Node
}

func (sensor *Sensor) canPropagateSignal() bool {
	return len(sensor.inbound) == 1 
}

func (sensor *Sensor) propagateSignal() {

	// this implemenation is just a stub which makes it easy to test for now.
	// at some point, sensors will act as proxies to real virtual sensors,
	// and probably be reading their inputs from sockets.

	if value, ok := <- sensor.inbound[0].channel; ok {
		sensor.scatterOutput(value)
	} 



}

// implementation needed here because when it was a method on *Node, it was calling 
// connectInboundWithChannel with a *Node instance and losing the fact it was a sensor
func (sensor *Sensor) ConnectBidirectional(target Connector) {
	sensor.ConnectBidirectionalWeighted(target, nil)
}

// implementation needed here because when it was a method on *Node, it was calling 
// connectInboundWithChannel with a *Node instance and losing the fact it was a sensor
func (sensor *Sensor) ConnectBidirectionalWeighted(target Connector, weights []float64) {
	channel := make(VectorChannel)		
	sensor.connectOutboundWithChannel(target, channel)
	target.connectInboundWithChannel(sensor, channel, weights)
}

