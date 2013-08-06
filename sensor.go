package neurgo

import (
	"fmt"
	"log"
)

// TODO: need a "function" which is called to gather actual data

type SensorFunction func(int) []float64

type Sensor struct {
	NodeId         *NodeId
	Outbound       []*OutboundConnection
	VectorLength   uint
	Closing        chan chan bool
	SyncChan       chan bool
	SensorFunction SensorFunction
}

func (sensor *Sensor) Run() {
	sensor.checkRunnable()

	closed := false
	syncCounter := 0

	for {

		select {
		case responseChan := <-sensor.Closing:
			closed = true
			responseChan <- true
			break // TODO: do we need this for anything??
		case _ = <-sensor.SyncChan:
			input := sensor.SensorFunction(syncCounter)
			syncCounter += 1
			dataMessage := &DataMessage{
				SenderId: sensor.NodeId,
				Inputs:   input,
			}
			sensor.scatterOutput(dataMessage)
		}

		if closed {
			sensor.Closing = nil
			sensor.SyncChan = nil
			break
		}
	}

}

func (sensor *Sensor) Init() {
	if sensor.Closing == nil {
		sensor.Closing = make(chan chan bool)
	} else {
		msg := "Warn: %v Init() called, but already had closing channel"
		log.Printf(msg, sensor)
	}

	if sensor.SyncChan == nil {
		sensor.SyncChan = make(chan bool)
	} else {
		msg := "Warn: %v Init() called, but already had data channel"
		log.Printf(msg, sensor)
	}
}

func (sensor *Sensor) checkRunnable() {
	if sensor.NodeId == nil {
		msg := fmt.Sprintf("not expecting sensor.NodeId to be nil")
		panic(msg)
	}

	if sensor.Closing == nil {
		msg := fmt.Sprintf("not expecting sensor.Closing to be nil")
		panic(msg)
	}

	if sensor.SyncChan == nil {
		msg := fmt.Sprintf("not expecting sensor.SyncChan to be nil")
		panic(msg)
	}

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

func (sensor *Sensor) scatterOutput(dataMessage *DataMessage) {
	for _, outboundConnection := range sensor.Outbound {
		dataChan := outboundConnection.DataChan
		log.Printf("Sensor %v scatter %v to: %v", sensor, dataMessage, outboundConnection)
		dataChan <- dataMessage
	}
}

func (sensor *Sensor) nodeId() *NodeId {
	return sensor.NodeId
}
