package neurgo

import (
	"sync"
	"fmt"
)

type NeuralNetwork struct {
	sensors   []*Node
	actuators []*Node
	Node
}

type copyScaffold struct {
	nodeScaffold map[*Node]*Node
	channelScaffold map[VectorChannel]VectorChannel
}


/*func (wiretap *Wiretap) propagateSignal() {}
func (wiretap *Wiretap) canPropagateSignal() bool { return false }
func (injector *Injector) propagateSignal() {}
func (injector *Injector) canPropagateSignal() bool { return false }*/


// Make sure the neural network gives expected output for the given 
// training samples.
func (neuralNet *NeuralNetwork) Verify(samples []*TrainingSample) bool {

	// make as many injectors as there are sensors
	injectors := make([]*Node, len(neuralNet.sensors))
	for i, _ := range injectors {
		injectors[i] = &Node{}
		injectors[i].Name = fmt.Sprintf("injector-%d", i+1)
		injectors[i].ConnectBidirectional(neuralNet.sensors[i])
	}

	// make as many wiretaps as actuators
	wiretaps := make([]*Node, len(neuralNet.actuators))
	for i, _ := range wiretaps {
		wiretaps[i] = &Node{}
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

func (neuralNet *NeuralNetwork) Run() {

	// get list of unique nodes in network
	nodes := neuralNet.uniqueNodes()

	// call Run() on each node
	for _, node := range nodes {
		Run(node.processor, node)
	}
	
	// NOTE: could refactor into call which performs a function on 
	// every unique node in the network 

}

func (neuralNet *NeuralNetwork) uniqueNodes() []*Node {
	return []*Node{neuralNet.sensors[0]}

}

func (neuralNet *NeuralNetwork) Copy() *NeuralNetwork {

	
	// the source neural network provides a "scaffold" for building the 
	// target network.  these provide the mapping between nodes and channels.
	nodeScaffold := make(map[*Node]*Node)
	channelScaffold := make(map[VectorChannel]VectorChannel)

	scaffold := &copyScaffold{nodeScaffold: nodeScaffold, channelScaffold: channelScaffold}

	sensorsCopy := make([]*Node,0)
	actuatorsCopy := make([]*Node, 0)
	neuralNetCopy := &NeuralNetwork{sensors: sensorsCopy, actuators: actuatorsCopy}

	for _, sensor := range neuralNet.sensors {
		sensorCopy := new(Node)
		// TODO: copy processor
		nodeScaffold[sensor] = sensorCopy
		sensorCopy.Name = sensor.Name
		neuralNetCopy.sensors = append(neuralNetCopy.sensors, sensorCopy)
	}

	for _, actuator := range neuralNet.actuators {
		actuatorCopy := new(Node)
		// TODO: copy processor
		nodeScaffold[actuator] = actuatorCopy
		actuatorCopy.Name = actuator.Name
		neuralNetCopy.actuators = append(neuralNetCopy.actuators, actuatorCopy)
	}

	for _, sensor := range neuralNet.sensors {
		sensorCopy := nodeScaffold[sensor]
		recreateOutboundConnectionsRecursive(sensor, sensorCopy, scaffold)
	}

	for _, actuator := range neuralNet.actuators {
		actuatorCopy := nodeScaffold[actuator]
		recreateInboundConnectionsRecursive(actuator, actuatorCopy, scaffold)
	}

	return neuralNetCopy

}


func recreateInboundConnectionsRecursive(nodeOriginal *Node, nodeCopy *Node, scaffold *copyScaffold) {
	
	nodeScaffold := scaffold.nodeScaffold
	channelScaffold := scaffold.channelScaffold

	for _, inboundConnection := range nodeOriginal.inboundConnections() {

		cxnTargetOriginal := inboundConnection.other
		cxnTargetCopy := createConnectionTargetCopy(cxnTargetOriginal, nodeScaffold)

		newCxn := &connection{}
		newCxn.other = cxnTargetCopy

		channelCopy := createChannelCopy(inboundConnection.channel, channelScaffold)
		newCxn.channel = channelCopy

		if inboundConnection.weights != nil && len(inboundConnection.weights) > 0 {
			weightsCopy := make([]float64, len(inboundConnection.weights))
			copy(weightsCopy, inboundConnection.weights)
			newCxn.weights = weightsCopy
		}
		
		nodeCopy.appendInboundConnection(newCxn)

		if len(cxnTargetOriginal.inboundConnections()) > 0 {
			recreateInboundConnectionsRecursive(cxnTargetOriginal, cxnTargetCopy, scaffold)
		} 
	} 
}


func recreateOutboundConnectionsRecursive(nodeOriginal *Node, nodeCopy *Node, scaffold *copyScaffold) {
	

	nodeScaffold := scaffold.nodeScaffold
	channelScaffold := scaffold.channelScaffold

	for _, outboundConnection := range nodeOriginal.outboundConnections() {

		cxnTargetOriginal := outboundConnection.other
		cxnTargetCopy := createConnectionTargetCopy(cxnTargetOriginal, nodeScaffold)

		newCxn := &connection{}
		newCxn.other = cxnTargetCopy

		channelCopy := createChannelCopy(outboundConnection.channel, channelScaffold)
		newCxn.channel = channelCopy

		nodeCopy.appendOutboundConnection(newCxn)

		if len(cxnTargetOriginal.outboundConnections()) > 0 {
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


func createConnectionTargetCopy(cxnTargetOriginal *Node, nodeScaffold map[*Node]*Node) *Node {

	var cxnTargetCopy *Node
	if cxnTargetCopyTemp, ok := nodeScaffold[cxnTargetOriginal]; ok {  // TODO: hack
		cxnTargetCopy = cxnTargetCopyTemp
	} else {

		// the connection target does not exist in nodeScaffold, create it
		node := &Node{}
		node.Name = cxnTargetOriginal.Name
		cxnTargetCopy = node
		nodeScaffold[cxnTargetOriginal] = cxnTargetCopy
		
	}

	return cxnTargetCopy


}
