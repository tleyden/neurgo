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

	
	// create an empty map between old nodes and new nodes
	// var nodeScaffold map[*Node]*Node
	nodeScaffold := make(map[Connector]Connector)

	sensorsCopy := make([]*Sensor,0)
	actuatorsCopy := make([]*Actuator, 0)
	neuralNetCopy := &NeuralNetwork{sensors: sensorsCopy, actuators: actuatorsCopy}

	for _, sensor := range neuralNet.sensors {
		sensorCopy := new(Sensor)
		nodeScaffold[sensor] = sensorCopy
		sensorCopy.Name = sensor.Name
		neuralNetCopy.sensors = append(neuralNetCopy.sensors, sensorCopy)
		log.Printf("sensorOriginal: %v", sensor)
		log.Printf("sensorCopy: %v", sensorCopy)

		recreateOutboundConnectionsRecursive(sensor, sensorCopy, nodeScaffold)
	}

	return neuralNetCopy

}


func recreateOutboundConnectionsRecursive(nodeOriginal Connector, nodeCopy Connector, nodeScaffold map[Connector]Connector) {
	
	// for each outbound connection:
	//   see if connection target node copy exists in map, if not, create it
	//   duplicate the outgoing connection
	//   maybe: duplicate the incoming connection on the other side of this outbound 
	//          (or maybe this makes sense to do in a separate pass)
	//   recursive call for the connection target
	
	var cxnTargetCopy Connector
	for _, outboundConnection := range nodeOriginal.outboundConnections() {

		cxnTargetOriginal := outboundConnection.other
		
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
			default:

				msg := fmt.Sprintf("unexpected cxnTargetOriginal type: %T", t) 
				panic(msg)
			}
			nodeScaffold[cxnTargetOriginal] = cxnTargetCopy

		}

		log.Printf("map: %v", nodeScaffold)

		newConnection := &connection{}
		newConnection.other = cxnTargetCopy

		if len(cxnTargetOriginal.outboundConnections()) > 0 {
			log.Printf("recursing into recreateOutboundConnectionsRecursive with: %v, %v, %v ", cxnTargetOriginal, cxnTargetCopy, nodeScaffold)
			recreateOutboundConnectionsRecursive(cxnTargetOriginal, cxnTargetCopy, nodeScaffold)
		} else {
			log.Printf("not recursing, cxnTargetOriginal has no outbound connections")
		}
		

	} 



}
