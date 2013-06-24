
# Neurgo

[![Build Status](https://drone.io/github.com/tleyden/neurgo/status.png)](https://drone.io/github.com/tleyden/neurgo/latest)

A library for constructing Neural Networks in [Go](http://golang.org/)

![architecture_diagram.png](http://cl.ly/image/143P2G2i3i1a/neurgo.png)


# Project Goals:

* Feature parity with [DXNN2](https://github.com/CorticalComputer/DXNN2) (a Topology & Parameter Evolving Universal Learning Network in Erlang)
* Support traditional Backpropagation learning methods, in addition to Evolutionary based methods
* 100% test coverage
* Thorough documentation & examples

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

# Status

* Feedforward operation is complete
* Learning algorithms (Stochastic Hill Climbing) in progress.
* Stay tuned!

# Understanding the codebase

Start with TestNetwork in neural_network_test.go

# Related Work

[DXNN2](https://github.com/CorticalComputer/DXNN2) - Pure Erlang TPEULN (Topology & Parameter Evolving Universal Learning Network).  

# Related Publications

[Handbook of Neuroevolution Through Erlang](http://www.amazon.com/Handbook-Neuroevolution-Through-Erlang-Gene/dp/1461444624) _by Gene Sher_.
