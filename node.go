package neurgo

import (
	"encoding/json"
	"log"
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

	if node.processor == nil {
		log.Panicf("%v does not have a signal processor", node)
	}

	for {

		if node.hasBeenShutdown() {
			break
		}

		if node.processor.canPropagate(node) == false {
			log.Panicf("%v cannot propagate any signals", node)
		} else {
			isShutdown := node.processor.propagateSignal(node)
			if isShutdown {
				break
			}
		}

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

func (node *Node) disconnectOutbound(target *Node) {
	for i, connection := range node.outbound {
		if connection.other == target {
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
