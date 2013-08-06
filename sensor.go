package neurgo

import (
	"log"
)

// TODO: need a "function" which is called to gather actual data

type Sensor struct {
	NodeId       *NodeId
	Outbound     []*OutboundConnection
	VectorLength uint
}

func (s *Sensor) ConnectOutbound(connectable OutboundConnectable) {
	if s.Outbound == nil {
		s.Outbound = make([]*OutboundConnection, 0)
	}

	if connectable.dataChan() == nil {
		log.Panicf("Cannot make outbound connection, dataChan == nil")
	}

	connection := &OutboundConnection{
		NodeId:   connectable.nodeId(),
		DataChan: connectable.dataChan(),
	}

	s.Outbound = append(s.Outbound, connection)
}

func (sensor *Sensor) nodeId() *NodeId {
	return sensor.NodeId
}
