package neurgo

import (
	"encoding/json"
	"time"
)

type Node struct {
	Name      string
	inbound   []*connection
	outbound  []*connection
	processor SignalProcessor
}

func (node *Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Type      string          `json:"type"`
			Name      string          `json:"name"`
			Outbound  []*connection   `json:"outbound"`
			Inbound   []*connection   `json:"inbound"`
			Processor SignalProcessor `json:"processor"`
		}{
			Type:      "Node",
			Name:      node.Name,
			Outbound:  node.outbound,
			Inbound:   node.inbound,
			Processor: node.processor,
		})
}

// continually propagate incoming signals -> outgoing signals
func (node *Node) Run() {
	for {
		if !node.processor.canPropagateSignal(node) {
			time.Sleep(time.Second) // <-- watch advanced concurrency talk and remove
		} else {
			node.processor.propagateSignal(node)
		}
	}
}

func (node *Node) Shutdown() {

}

func (node *Node) String() string {
	return node.Name
}

func (node *Node) scatterOutput(outputs []float64) {
	for _, outboundConnection := range node.outbound {
		select {
		case outboundConnection.channel <- outputs:
		case <-outboundConnection.closing:
			return
		}

	}
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
	closing := make(chan bool)
	node.connectOutboundWithChannel(target, channel, closing)
	target.connectInboundWithChannel(node, channel, closing, weights)
}

func (node *Node) connectOutboundWithChannel(target *Node, channel VectorChannel, closing chan bool) {
	connection := &connection{
		channel: channel,
		other:   target,
		closing: closing,
	}
	node.outbound = append(node.outbound, connection)
}

func (node *Node) connectInboundWithChannel(source *Node, channel VectorChannel, closing chan bool, weights []float64) {
	connection := &connection{
		channel: channel,
		weights: weights,
		other:   source,
		closing: closing,
	}
	node.inbound = append(node.inbound, connection)
}

func (node *Node) DisconnectBidirectional(target *Node) {
	node.disconnectOutbound(target)
	target.disconnectInbound(node)
}

func (node *Node) disconnectOutbound(target *Node) {
	for i, connection := range node.outbound {
		if connection.other == target {
			close(node.outbound[i].closing)
			node.outbound = removeConnection(node.outbound, i)
		}
	}
}

func (node *Node) disconnectInbound(source *Node) {
	for i, connection := range node.inbound {
		if connection.other == source {
			node.inbound = removeConnection(node.inbound, i)
		}
	}
}

func removeConnection(connections []*connection, index int) []*connection {

	newConnections := make([]*connection, len(connections)-1)
	newConnectionsIndex := 0

	for i, connection := range connections {
		if i != index {
			newConnections[newConnectionsIndex] = connection
			newConnectionsIndex += 1
		}
	}
	return newConnections

}

func (node *Node) outboundConnections() []*connection {
	return node.outbound
}

func (node *Node) inboundConnections() []*connection {
	return node.inbound
}

func (node *Node) appendOutboundConnection(target *connection) {
	node.outbound = append(node.outbound, target)
}

func (node *Node) appendInboundConnection(source *connection) {
	node.inbound = append(node.inbound, source)
}
