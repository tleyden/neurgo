package neurgo

import (
	"sync"
	"fmt"
	"log"
)

type NeuralNetwork struct {
	sensors   []*Sensor
	actuators []*Actuator
	Node
}

type copyScaffold struct {
	nodeScaffold map[Connector]Connector
	channelScaffold map[VectorChannel]VectorChannel
}

type Wiretap struct {
	Node
}

type Injector struct {
	Node
}

// Make sure the neural network gives expected output for the given 
// training samples.
func (neuralNet *NeuralNetwork) Verify(samples []*TrainingSample) bool {

	// make as many injectors as there are sensors
	injectors := make([]*Injector, len(neuralNet.sensors))
	for i, _ := range injectors {
		injectors[i] = &Injector{}
		injectors[i].Name = fmt.Sprintf("injector-%d", i+1)
		injectors[i].ConnectBidirectional(neuralNet.sensors[i])
	}

	// make as many wiretaps as actuators
	wiretaps := make([]*Wiretap, len(neuralNet.actuators))
	for i, _ := range wiretaps {
		wiretaps[i] = &Wiretap{}
		wiretaps[i].Name = fmt.Sprintf("wiretap-%d", i+1)
		neuralNet.actuators[i].ConnectBidirectional(wiretaps[i])
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)

	// inject values into sensors
	go func() {
		for _, sample := range samples {
			for j, inputsForSensor := range sample.sampleInputs {
				injectors[j].outbound[0].channel <- inputsForSensor
			}
		}
		wg.Done()
	}()

	// read the value from wiretap (which taps into actuator)
	verified := true
	go func() {

		for _, sample := range samples {
			for j, expectedOutputs := range sample.expectedOutputs {
				resultVector := <- wiretaps[j].inbound[0].channel
				if !vectorEqualsWithMaxDelta(resultVector, expectedOutputs, 0.01) {
					verified = false
				}
			}
		}

		wg.Done() 
	}()

	wg.Wait()
	
	// disconnect injectors and wiretaps to leave it in the same state!
	for i, injector := range injectors {
		injector.DisconnectBidirectional(neuralNet.sensors[i])
	}
	for i, actuator := range neuralNet.actuators {
		actuator.DisconnectBidirectional(wiretaps[i])
	}

	return verified  

}

func (neuralNet *NeuralNetwork) Copy() *NeuralNetwork {

	
	// the source neural network provides a "scaffold" for building the 
	// target network.  these provide the mapping between nodes and channels.
	nodeScaffold := make(map[Connector]Connector)
	channelScaffold := make(map[VectorChannel]VectorChannel)

	scaffold := &copyScaffold{nodeScaffold: nodeScaffold, channelScaffold: channelScaffold}

	sensorsCopy := make([]*Sensor,0)
	actuatorsCopy := make([]*Actuator, 0)
	neuralNetCopy := &NeuralNetwork{sensors: sensorsCopy, actuators: actuatorsCopy}

	for _, sensor := range neuralNet.sensors {
		sensorCopy := new(Sensor)
		nodeScaffold[sensor] = sensorCopy
		sensorCopy.Name = sensor.Name
		neuralNetCopy.sensors = append(neuralNetCopy.sensors, sensorCopy)
	}

	for _, actuator := range neuralNet.actuators {
		actuatorCopy := new(Actuator)
		nodeScaffold[actuator] = actuatorCopy
		actuatorCopy.Name = actuator.Name
		neuralNetCopy.actuators = append(neuralNetCopy.actuators, actuatorCopy)
	}


	for _, sensor := range neuralNet.sensors {
		sensorCopy := nodeScaffold[sensor]
		recreateOutboundConnectionsRecursive(sensor, sensorCopy, scaffold)
	}

	log.Printf("-----------------------------------")

	// TODO: inbound connections!!
	for _, actuator := range neuralNet.actuators {
		actuatorCopy := nodeScaffold[actuator]
		recreateInboundConnectionsRecursive(actuator, actuatorCopy, scaffold)
	}

	return neuralNetCopy

}


func recreateInboundConnectionsRecursive(nodeOriginal Connector, nodeCopy Connector, scaffold *copyScaffold) {
	
	log.Printf("recreateInboundConnectionsRecursive called with: %v", nodeOriginal)
	nodeScaffold := scaffold.nodeScaffold
	channelScaffold := scaffold.channelScaffold

	for _, inboundConnection := range nodeOriginal.inboundConnections() {

		cxnTargetOriginal := inboundConnection.other
		cxnTargetCopy := createConnectionTargetCopy(cxnTargetOriginal, nodeScaffold)

		newCxn := &connection{}
		newCxn.other = cxnTargetCopy

		channelCopy := createChannelCopy(inboundConnection.channel, channelScaffold)
		newCxn.channel = channelCopy

		log.Printf("append inbound connection %v to %v", newCxn, nodeCopy)
		nodeCopy.appendInboundConnection(newCxn)
		log.Printf("%v (%p) now has %v inbound connections", nodeCopy, nodeCopy, len(nodeCopy.inboundConnections()))


		if len(cxnTargetOriginal.inboundConnections()) > 0 {
			log.Printf("recursing into recreateInboundConnectionsRecursive with: %v, %v, %v ", cxnTargetOriginal, cxnTargetCopy, nodeScaffold)
			recreateInboundConnectionsRecursive(cxnTargetOriginal, cxnTargetCopy, scaffold)
		} 
		

	} 

}


func recreateOutboundConnectionsRecursive(nodeOriginal Connector, nodeCopy Connector, scaffold *copyScaffold) {
	
	log.Printf("recreateOutboundConnectionsRecursive called with: %v", nodeOriginal)
	nodeScaffold := scaffold.nodeScaffold
	channelScaffold := scaffold.channelScaffold

	for _, outboundConnection := range nodeOriginal.outboundConnections() {

		cxnTargetOriginal := outboundConnection.other
		cxnTargetCopy := createConnectionTargetCopy(cxnTargetOriginal, nodeScaffold)

		newCxn := &connection{}
		newCxn.other = cxnTargetCopy

		channelCopy := createChannelCopy(outboundConnection.channel, channelScaffold)
		newCxn.channel = channelCopy

		log.Printf("append connection %v to %v", newCxn, nodeCopy)
		nodeCopy.appendOutboundConnection(newCxn)
		log.Printf("nodeCopy now has %v outbound connections", len(nodeCopy.outboundConnections()))


		if len(cxnTargetOriginal.outboundConnections()) > 0 {
			log.Printf("recursing into recreateOutboundConnectionsRecursive with: %v, %v, %v ", cxnTargetOriginal, cxnTargetCopy, nodeScaffold)
			recreateOutboundConnectionsRecursive(cxnTargetOriginal, cxnTargetCopy, scaffold)
		} 
		

	} 

}

func createChannelCopy(channelOriginal VectorChannel, channelScaffold map[VectorChannel]VectorChannel) VectorChannel {

	var channelCopy VectorChannel
	if channelCopyTemp, ok := channelScaffold[channelOriginal]; ok {
		channelCopy = channelCopyTemp
	} else {
		channelCopy = make(VectorChannel)
		channelScaffold[channelOriginal] = channelCopy
	}
	return channelCopy

}


func createConnectionTargetCopy(cxnTargetOriginal Connector, nodeScaffold map[Connector]Connector) Connector {

	var cxnTargetCopy Connector
	if cxnTargetCopyTemp, ok := nodeScaffold[cxnTargetOriginal]; ok {  // TODO: hack
		cxnTargetCopy = cxnTargetCopyTemp
	} else {

		// the connection target does not exist in nodeScaffold, create it
		switch t:= cxnTargetOriginal.(type) {
		case *Sensor:
			log.Printf("its a sensor: %T", t)
			sensor := &Sensor{}
			sensor.Name = t.Name  
			cxnTargetCopy = sensor
		case *Neuron:
			log.Printf("its a neuron: %T", t)
			neuron := &Neuron{}
			neuron.Name = t.Name
			cxnTargetCopy = neuron
		case *Actuator:
			log.Printf("its an actuator: %T", t)
			actuator := &Actuator{}
			actuator.Name = t.Name
			cxnTargetCopy = actuator
		case *Node:
			log.Printf("its a node: %T %v", t, t)
			node := &Node{}
			node.Name = t.Name
			cxnTargetCopy = node

		default:
			msg := fmt.Sprintf("unexpected cxnTargetOriginal type: %T", t) 
			panic(msg)
		}
		nodeScaffold[cxnTargetOriginal] = cxnTargetCopy
		
	}

	return cxnTargetCopy


}
