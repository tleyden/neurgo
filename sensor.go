package neurgo

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type SensorFunction func(int) []float64

type Sensor struct {
	NodeId         *NodeId
	Outbound       []*OutboundConnection
	VectorLength   uint
	Closing        chan chan bool
	SyncChan       chan bool
	SensorFunction SensorFunction
	wg             sync.WaitGroup
}

func (sensor *Sensor) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			NodeId       *NodeId
			VectorLength uint
			Outbound     []*OutboundConnection
		}{
			NodeId:       sensor.NodeId,
			VectorLength: sensor.VectorLength,
			Outbound:     sensor.Outbound,
		})
}

func (sensor *Sensor) String() string {
	return JsonString(sensor)
}

func (sensor *Sensor) Run() {

	defer sensor.wg.Done()

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
			log.Printf("Sensor got value from SyncChan")
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

	if sensor.SensorFunction == nil {
		// if there is no SensorFunction, create a default
		// function which emits a 0-vector
		sensorFunc := func(syncCounter int) []float64 {
			return make([]float64, sensor.VectorLength)
		}
		sensor.SensorFunction = sensorFunc
	}

	sensor.wg.Add(1) // TODO: make sure Init() not called twice!
}

func (sensor *Sensor) Shutdown() {

	closingResponse := make(chan bool)
	sensor.Closing <- closingResponse
	response := <-closingResponse
	if response != true {
		log.Panicf("Got unexpected response on closing channel")
	}

	sensor.wg.Wait()
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

	if sensor.SensorFunction == nil {
		msg := fmt.Sprintf("not expecting sensor.SensorFunction to be nil")
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
