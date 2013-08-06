package neurgo

import (
	"log"
)

// TODO: need a "function" which is called to actuate based on data

type Actuator struct {
	NodeId       *NodeId
	Inbound      []*InboundConnection
	Closing      chan chan bool
	DataChan     chan *DataMessage
	VectorLength uint
}

func (actuator *Actuator) Init() {
	if actuator.Closing == nil {
		actuator.Closing = make(chan chan bool)
	} else {
		msg := "Warn: %v Init() called, but already had closing channel"
		log.Printf(msg, actuator)
	}

	if actuator.DataChan == nil {
		actuator.DataChan = make(chan *DataMessage, len(actuator.Inbound))
	} else {
		msg := "Warn: %v Init() called, but already had data channel"
		log.Printf(msg, actuator)
	}
}

func (actuator *Actuator) ConnectInbound(connectable InboundConnectable) {
	if actuator.Inbound == nil {
		actuator.Inbound = make([]*InboundConnection, 0)
	}
	connection := &InboundConnection{
		NodeId:  connectable.nodeId(),
		Weights: nil,
	}
	actuator.Inbound = append(actuator.Inbound, connection)
}

func (actuator *Actuator) dataChan() chan *DataMessage {
	return actuator.DataChan
}

func (actuator *Actuator) nodeId() *NodeId {
	return actuator.NodeId
}
