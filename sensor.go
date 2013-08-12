package neurgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

type SensorFunction func(int) []float64

type Sensor struct {
	NodeId         *NodeId
	Outbound       []*OutboundConnection
	VectorLength   int
	Closing        chan chan bool
	SyncChan       chan bool
	SensorFunction SensorFunction
	wg             *sync.WaitGroup
	Cortex         *Cortex
}

func (sensor *Sensor) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			NodeId       *NodeId
			VectorLength int
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

	log.Printf("%v Run() started", sensor.NodeId.UUID)

	defer sensor.wg.Done()

	sensor.checkRunnable()

	closed := false
	syncCounter := 0

	for {
		select {
		case responseChan := <-sensor.Closing:
			log.Printf("%v received Closing message", sensor.NodeId.UUID)
			closed = true
			responseChan <- true
			break // TODO: do we need this for anything??
		case _ = <-sensor.SyncChan:
			log.Printf("%v received Sync message", sensor.NodeId.UUID)
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

	log.Printf("%v Run() finished", sensor.NodeId.UUID)

}

func (sensor *Sensor) Init(reInit bool) {
	if reInit == true {
		sensor.Closing = make(chan chan bool)
	} else if sensor.Closing == nil {
		sensor.Closing = make(chan chan bool)
	}

	if reInit == true {
		sensor.SyncChan = make(chan bool)
	} else {
		if sensor.SyncChan == nil {
			sensor.SyncChan = make(chan bool)
		}
	}

	if reInit == true {
		sensorFunc := func(syncCounter int) []float64 {
			return make([]float64, sensor.VectorLength)
		}
		sensor.SensorFunction = sensorFunc

	} else {
		if sensor.SensorFunction == nil {
			// if there is no SensorFunction, create a default
			// function which emits a 0-vector
			sensorFunc := func(syncCounter int) []float64 {
				return make([]float64, sensor.VectorLength)
			}
			sensor.SensorFunction = sensorFunc
		}
	}

	if reInit == true {
		sensor.wg = &sync.WaitGroup{}
		sensor.wg.Add(1)

	} else {
		if sensor.wg == nil {
			sensor.wg = &sync.WaitGroup{}
			sensor.wg.Add(1)
		}

	}

}

func (sensor *Sensor) Shutdown() {

	closingResponse := make(chan bool)
	sensor.Closing <- closingResponse
	response := <-closingResponse
	if response != true {
		log.Panicf("Got unexpected response on closing channel")
	}

	sensor.shutdownOutboundConnections()

	sensor.wg.Wait()
	sensor.wg = nil
}

func (sensor *Sensor) outbound() []*OutboundConnection {
	return sensor.Outbound
}

func (sensor *Sensor) setOutbound(newOutbound []*OutboundConnection) {
	sensor.Outbound = newOutbound
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

	if err := sensor.validateOutbound(); err != nil {
		msg := fmt.Sprintf("invalid outbound connection(s): %v", err.Error())
		panic(msg)
	}

}

func (sensor *Sensor) validateOutbound() error {
	for _, connection := range sensor.Outbound {
		if connection.DataChan == nil {
			msg := fmt.Sprintf("%v has empty DataChan", connection)
			return errors.New(msg)
		}
	}
	return nil
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
	log.Printf("Sensor scatterOutbound to: %v", sensor.Outbound)
	for _, outboundConnection := range sensor.Outbound {
		dataChan := outboundConnection.DataChan
		log.Printf("Sensor %v start scatter %v to: %v dataChan: %v", sensor.NodeId.UUID, dataMessage, outboundConnection, dataChan)
		dataChan <- dataMessage
		log.Printf("Sensor %v finish scatter %v to: %v dataChan: %v", sensor.NodeId.UUID, dataMessage, outboundConnection, dataChan)
	}
}

func (sensor *Sensor) nodeId() *NodeId {
	return sensor.NodeId
}

func (sensor *Sensor) initOutboundConnections(nodeIdToDataMsg nodeIdToDataMsgMap) {
	for _, outboundConnection := range sensor.Outbound {
		if outboundConnection.DataChan == nil {
			dataChan := nodeIdToDataMsg[outboundConnection.NodeId.UUID]
			if dataChan != nil {
				outboundConnection.DataChan = dataChan
			}
		}
	}
}

func (sensor *Sensor) shutdownOutboundConnections() {
	for _, outboundConnection := range sensor.Outbound {
		outboundConnection.DataChan = nil
	}
}
