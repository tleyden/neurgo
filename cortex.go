package neurgo

type Neuron struct {
	Sensors   []NodeId
	Neurons   []NodeId
	Actuators []NodeId
}
