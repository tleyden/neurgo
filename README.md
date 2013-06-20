
# Neurgo

A library for constructing Neural Networks in [Go](http://golang.org/)

![architecture_diagram.png](http://cl.ly/image/143P2G2i3i1a/neurgo.png)


# Project Goals:

* Feature parity with [DXNN2](https://github.com/CorticalComputer/DXNN2) (a topology and weight evolving neural network in Erlang)
* Support traditional Backpropagation learning methods, in addition to Evolutionary based methods
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

# Status

* Feedforward operation is complete
* Learning algorithms (Stochastic Hill Climbing) in progress.
* Stay tuned!

# Understanding the codebase

Start with TestNetwork in neural_network_test.go

# Related Work

[DXNN2](https://github.com/CorticalComputer/DXNN2) - Pure Erlang TWEANN (Topology WEight Adapting Neural Network).  

# Related Publications

[Handbook of Neuroevolution Through Erlang](http://www.amazon.com/Handbook-Neuroevolution-Through-Erlang-Gene/dp/1461444624) _by Gene Sher_.
