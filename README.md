
# Neurgo

A library for constructing Neural Networks in [Go](http://golang.org/)

The end goal:

* Feature parity with [DXNN2](https://github.com/CorticalComputer/DXNN2) (a topology and weight evolving neural network in Erlang)
* Support traditional Backpropagation learning methods
* 100% test coverage
* Thorough documentation & examples

# Example code

```
// create network nodes
neuron1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
neuron2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
sensor := &Sensor{}
actuator := &Actuator{}

// connect nodes together
weights := []float64{20,20,20,20,20}
sensor.ConnectBidirectionalWeighted(neuron1, weights)
sensor.ConnectBidirectionalWeighted(neuron2, weights)
neuron1.ConnectBidirectional(actuator)
neuron2.ConnectBidirectional(actuator)

// spinup node goroutines
go Run(neuron1)
go Run(neuron2)
go Run(sensor)
go Run(actuator)
```

# Architecture

* Each node in the network is a go-routine
* Nodes communicate with eachother over channels

# Status

* Feedforward operation is nearly complete
* No learning algorithms have been added yet, so it's basically useless (braindead)
* Stay tuned!

# Understanding the codebase

Start with TestNetwork in node_test.go
