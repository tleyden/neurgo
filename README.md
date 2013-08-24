# Neurgo

[![Build Status](https://drone.io/github.com/tleyden/neurgo/status.png)](https://drone.io/github.com/tleyden/neurgo/latest)

A library for constructing Neural Networks in [Go](http://golang.org/)

![architecture_diagram.png](http://cl.ly/image/143P2G2i3i1a/neurgo.png)

For a more detailed architecture diagram, see the [Neurgo Architecture Prezi](http://prezi.com/cldumvoxwsxj/?utm_campaign=share&utm_medium=copy)

# Project Goals:

* Feature parity with [DXNN2](https://github.com/CorticalComputer/DXNN2) (a Topology & Parameter Evolving Universal Learning Network in Erlang)
* Support traditional Backpropagation learning methods, in addition to Evolutionary based methods
* 100% test coverage
* Message passing architecture 
* Complete documentation & examples

# Example code

```
// create network nodes
neuronProcessor1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
neuronProcessor2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
neuron1 := &Node{Name: "neuron1", processor: neuronProcessor1}
neuron2 := &Node{Name: "neuron2", processor: neuronProcessor2}
sensor := &Node{Name: "sensor", processor: &Sensor{}}
actuator := &Node{Name: "actuator", processor: &Actuator{}}

// connect nodes together
weights := []float64{20, 20, 20, 20, 20}
sensor.ConnectBidirectionalWeighted(neuron1, weights)
sensor.ConnectBidirectionalWeighted(neuron2, weights)
neuron1.ConnectBidirectional(actuator)
neuron2.ConnectBidirectional(actuator)

// create neural network
sensors := []*Node{sensor}
actuators := []*Node{actuator}
neuralNet := &NeuralNetwork{sensors: sensors, actuators: actuators}

// spin up node goroutines
neuralNet.Run()

// inputs + expected outputs
examples := []*TrainingSample{{sampleInputs: [][]float64{[]float64{1, 1, 1, 1, 1}}, expectedOutputs: [][]float64{[]float64{110, 110}}}}

// verify neural network
verified := neuralNet.Verify(examples)
assert.True(t, verified)
        
```

# Getting Started

* [Install Go](http://golang.org/doc/install)

* Clone repository with `$ git clone git://github.com/tleyden/neurgo.git`

* Run tests with `$ go test`

* To write code that uses neurgo, your code will need `import "github.com/tleyden/neurgo"` as described in the [API documentation](http://godoc.org/github.com/tleyden/neurgo)

# Documentation

* This README file

* [API documentation](http://godoc.org/github.com/tleyden/neurgo)


# Status

* Feedforward operation works
* Learning via Stochastic Hill Climbing works (see [neurvolve](https://github.com/tleyden/neurvolve))
* Topological mutatation operators implemented (see [neurvolve](https://github.com/tleyden/neurvolve))
* Example which evolves a network capable of sovling XOR (see [neurvolve](https://github.com/tleyden/neurvolve))

# TODO

* Visualize networks
* Backpropagation based learning (contributions welcome)

# Libraries that build on Neurgo

* [neurvolve](https://github.com/tleyden/neurvolve) builds on this library

# Related Work

[DXNN2](https://github.com/CorticalComputer/DXNN2) - Pure Erlang TPEULN (Topology & Parameter Evolving Universal Learning Network).  


# Related Publications

[Handbook of Neuroevolution Through Erlang](http://www.amazon.com/Handbook-Neuroevolution-Through-Erlang-Gene/dp/1461444624) _by Gene Sher_.


