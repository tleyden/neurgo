package neurgo

import (
	"encoding/json"
	"sync"
)

type Node struct {
	Name      string
	inbound   []*Connection
	outbound  []*Connection
	processor SignalProcessor
	closing   chan bool
	invisible bool
	wg        sync.WaitGroup
}

type weightedInput struct {
	weights []float64
	inputs  []float64
}

// continually propagate incoming signals -> outgoing signals
func (node *Node) Run() {

	// create the closing channel before goroutine kicked off
	// to solve a potential race condition where someone tries
	// to shutdown a not-yet-running node
	node.closing = make(chan bool)
	node.wg.Add(1)

	go node.runGoroutine()

}

func (node *Node) runGoroutine() {

	defer node.wg.Done()

	panicIfNil(node.processor)

	for {

		weightedInputs := make([]*weightedInput, 0)
		isShutdown := false

		for _, connection := range node.inbound {

			var inputs []float64
			select {
			case inputs = <-connection.channel:
				panicIfZero(len(inputs))
			case <-node.closing:
				isShutdown = true
				break
			}

			weights := connection.weights
			weightedInput := &weightedInput{weights: weights, inputs: inputs}
			weightedInputs = append(weightedInputs, weightedInput)

		}

		if isShutdown {
			break
		}

		outputs := node.processor.CalculateOutput(weightedInputs)
		node.scatterOutput(outputs)

	}

}

func (node *Node) Shutdown() {
	close(node.closing)
	node.wg.Wait()
}

func (node *Node) hasBeenShutdown() bool {
	select {
	case <-node.closing:
		return true
	default:
		return false
	}

}

func (node *Node) String() string {
	return node.Name
}

func (node *Node) scatterOutput(outputs []float64) {
	for _, outboundConnection := range node.outbound {
		select {
		case outboundConnection.channel <- outputs:
		case <-node.closing:
			return
		}

	}
}

func (node *Node) isInvisible() bool {
	return node.invisible
}

func (node *Node) setInvisible(val bool) {
	node.invisible = val
}

func (node *Node) hasOutboundConnectionTo(other *Node) bool {
	for _, outboundConnection := range node.outbound {
		if outboundConnection.other == other {
			return true
		}
	}
	return false
}

func (node *Node) ConnectBidirectional(target *Node) {
	node.ConnectBidirectionalWeighted(target, nil)
}

func (node *Node) ConnectBidirectionalWeighted(target *Node, weights []float64) {
	channel := make(VectorChannel)
	node.connectOutboundWithChannel(target, channel)
	target.connectInboundWithChannel(node, channel, weights)
}

func (node *Node) connectOutboundWithChannel(target *Node, channel VectorChannel) {
	connection := &Connection{
		channel: channel,
		other:   target,
	}
	node.outbound = append(node.outbound, connection)
}

func (node *Node) connectInboundWithChannel(source *Node, channel VectorChannel, weights []float64) {
	connection := &Connection{
		channel: channel,
		weights: weights,
		other:   source,
	}
	node.inbound = append(node.inbound, connection)
}

func (node *Node) DisconnectBidirectional(target *Node) {
	node.disconnectOutbound(target)
	target.disconnectInbound(node)
}

func (node *Node) DisconnectAllBidirectional() {
	for _, connection := range node.outbound {
		other := connection.other
		node.disconnectOutbound(other)
		other.disconnectInbound(node)
	}
}

func (node *Node) disconnectOutbound(target *Node) {
	for _, connection := range node.outbound {
		if connection.other == target {
			node.outbound = removeConnection(node.outbound, connection)
		}
	}
}

func (node *Node) disconnectInbound(source *Node) {

	// TODO: find connections to remove in one step
	// remove connections in another step (loop over each connection to remove)
	// change removeConnection around to take a connection instance instead of index

	for _, connection := range node.inbound {
		if connection.other == source {
			node.inbound = removeConnection(node.inbound, connection)
		}
	}
}

func removeConnection(connections []*Connection, removing *Connection) []*Connection {
	newConnections := make([]*Connection, 0)
	for _, connection := range connections {
		if connection != removing {
			newConnections = append(newConnections, connection)
		}
	}
	return newConnections
}

/*
func removeConnection(connections []*Connection, index int) []*Connection {

	newConnections := make([]*Connection, len(connections)-1)
	newConnectionsIndex := 0

	for i, connection := range connections {
		if i != index {
			newConnections[newConnectionsIndex] = connection
			newConnectionsIndex += 1
		}
	}
	return newConnections

}
*/

func (node *Node) outboundConnections() []*Connection {
	return node.outbound
}

func (node *Node) inboundConnections() []*Connection {
	return node.inbound
}

func (node *Node) appendOutboundConnection(target *Connection) {
	node.outbound = append(node.outbound, target)
}

func (node *Node) appendInboundConnection(source *Connection) {
	node.inbound = append(node.inbound, source)
}

func (node *Node) Inbound() []*Connection {
	return node.inbound
}

func (node *Node) Outbound() []*Connection {
	return node.outbound
}

func (node *Node) Processor() SignalProcessor {
	return node.processor
}

func (node *Node) SetProcessor(processor SignalProcessor) {
	node.processor = processor
}

func (node *Node) IsNeuron() bool {
	// only neurons have bias, so do a quick hack which
	// leverages that tidbit of knowledge.
	return node.processor.HasBias()
}

func (node *Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Type      string          `json:"type"`
			Name      string          `json:"name"`
			Outbound  []*Connection   `json:"outbound"`
			Inbound   []*Connection   `json:"inbound"`
			Processor SignalProcessor `json:"processor"`
		}{
			Type:      "Node",
			Name:      node.Name,
			Outbound:  node.outbound,
			Inbound:   node.inbound,
			Processor: node.processor,
		})
}
