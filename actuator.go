package neurgo

import (
	"fmt"
	"log"
)

type ActuatorFunction func(outputs []float64)

type Actuator struct {
	NodeId           *NodeId
	Inbound          []*InboundConnection
	Closing          chan chan bool
	DataChan         chan *DataMessage
	VectorLength     uint
	ActuatorFunction ActuatorFunction
}

func (actuator *Actuator) Run() {

	log.Printf("")

	actuator.checkRunnable()

	weightedInputs := createEmptyWeightedInputs(actuator.Inbound)

	closed := false

	for {

		log.Printf("Actuator %v select().  datachan: %v", actuator, actuator.DataChan)

		select {
		case responseChan := <-actuator.Closing:
			closed = true
			responseChan <- true
			break // TODO: do we need this for anything??
		case dataMessage := <-actuator.DataChan:
			recordInput(weightedInputs, dataMessage)
		}

		if closed {
			actuator.Closing = nil
			actuator.DataChan = nil
			break
		}

		if receiveBarrierSatisfied(weightedInputs) {

			log.Printf("Actuator %v received inputs: %v", actuator, weightedInputs)
			scalarOutput := actuator.computeScalarOutput(weightedInputs)

			actuator.ActuatorFunction(scalarOutput)

			weightedInputs = createEmptyWeightedInputs(actuator.Inbound)

		} else {
			log.Printf("Actuator %v receive barrier not satisfied", actuator)
		}

	}

}

func (actuator *Actuator) computeScalarOutput(weightedInputs []*weightedInput) []float64 {

	outputs := make([]float64, 0)
	for _, weightedInput := range weightedInputs {
		inputs := weightedInput.inputs
		actuator.validateInputs(inputs)
		inputValue := inputs[0]
		outputs = append(outputs, inputValue)
	}

	return outputs

}

func (actuator *Actuator) validateInputs(inputs []float64) {
	if len(inputs) != 1 {
		t := "%T got invalid input vector: %v"
		message := fmt.Sprintf(t, actuator, inputs)
		panic(message)
	}
}

func (actuator *Actuator) checkRunnable() {
	if actuator.NodeId == nil {
		msg := fmt.Sprintf("not expecting actuator.NodeId to be nil")
		panic(msg)
	}

	if actuator.Closing == nil {
		msg := fmt.Sprintf("not expecting actuator.Closing to be nil")
		panic(msg)
	}

	if actuator.DataChan == nil {
		msg := fmt.Sprintf("not expecting actuator.DataChan to be nil")
		panic(msg)
	}

	if actuator.ActuatorFunction == nil {
		msg := fmt.Sprintf("not expecting actuator.ActuatorFunction to be nil")
		panic(msg)
	}

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

func (actuator *Actuator) String() string {
	return fmt.Sprintf("%v", actuator.NodeId)
}
