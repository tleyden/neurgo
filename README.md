# Neurgo

[![Build Status](https://drone.io/github.com/tleyden/neurgo/status.png)](https://drone.io/github.com/tleyden/neurgo/latest) [![GoDoc](https://godoc.org/github.com/tleyden/neurgo?status.png)](https://godoc.org/github.com/tleyden/neurgo)

A library for constructing Neural Networks in [Go](http://golang.org/), where Neurons are goroutines that communicate with each other via channels.


![architecture_diagram.png](http://cl.ly/image/0a1Y0e0B0P1m/Screen%20Shot%202013-10-09%20at%209.22.11%20PM.png)

## What it can do

* Feedforward networks
* Recurrent networks
* JSON Marshal/Unmarshal ([example json](https://drone.io/github.com/tleyden/neurgo/files/xnor.json))
* Visualization network topology in SVG ([example svg](https://drone.io/github.com/tleyden/neurgo/files/xnor.svg))

## Learning mechanism

Neurgo does _not_ contain any code for learning/training.  

The idea is to have a separation of concerns such that the code that does the training will live in it's own repo.  Currently, there is only one training module:

* [neurvolve](https://github.com/tleyden/neurvolve) - An evolution based trainer that is essentially a port of [DXNN2](https://github.com/CorticalComputer/DXNN2) (a Topology & Parameter Evolving Universal Learning Network in Erlang).

## Roadmap

* Training module for Backpropagation based learning (contributions welcome!)
* Stress testing / benchmarks

## Example applications

* [Checkerlution - A Checkers Bot](https://github.com/tleyden/checkerlution)

## Example code

The following code creates a neural net with [this topology](https://drone.io/github.com/tleyden/neurgo/files/xnor.svg).  It does not actually run the network (eg, feed inputs), so for a more complete example see `cortex_test.go`.

```go
sensor := &Sensor{
	NodeId:       NewSensorId("sensor", 0.0),
	VectorLength: 2,
}
sensor.Init()
hiddenNeuron1 := &Neuron{
	ActivationFunction: EncodableSigmoid(),
	NodeId:             NewNeuronId("hidden-neuron1", 0.25),
	Bias:               -30,
}
hiddenNeuron1.Init()
hiddenNeuron2 := &Neuron{
	ActivationFunction: EncodableSigmoid(),
	NodeId:             NewNeuronId("hidden-neuron2", 0.25),
	Bias:               10,
}
hiddenNeuron2.Init()
outputNeuron := &Neuron{
	ActivationFunction: EncodableSigmoid(),
	NodeId:             NewNeuronId("output-neuron", 0.35),
	Bias:               -10,
}
outputNeuron.Init()
actuator := &Actuator{
	NodeId:       NewActuatorId("actuator", 0.5),
	VectorLength: 1,
}
actuator.Init()

// wire up connections
sensor.ConnectOutbound(hiddenNeuron1)
hiddenNeuron1.ConnectInboundWeighted(sensor, []float64{20, 20})
sensor.ConnectOutbound(hiddenNeuron2)
hiddenNeuron2.ConnectInboundWeighted(sensor, []float64{-20, -20})
hiddenNeuron1.ConnectOutbound(outputNeuron)
outputNeuron.ConnectInboundWeighted(hiddenNeuron1, []float64{20})
hiddenNeuron2.ConnectOutbound(outputNeuron)
outputNeuron.ConnectInboundWeighted(hiddenNeuron2, []float64{20})
outputNeuron.ConnectOutbound(actuator)
actuator.ConnectInbound(outputNeuron)

// create cortex
nodeId := NewCortexId("cortex")
cortex := &Cortex{
	NodeId: nodeId,
}
cortex.SetSensors([]*Sensor{sensor})
cortex.SetNeurons([]*Neuron{hiddenNeuron1, hiddenNeuron2, outputNeuron})
cortex.SetActuators([]*Actuator{actuator})
```

## Getting Started

* [Install Go](http://golang.org/doc/install)

* Clone repository with `$ git clone git://github.com/tleyden/neurgo.git`

* Run tests with `$ go test`

* To write code that uses neurgo, your code will need `import "github.com/tleyden/neurgo"` as described in the [API documentation](http://godoc.org/github.com/tleyden/neurgo)

## Documentation

* This README file

* [API documentation](http://godoc.org/github.com/tleyden/neurgo)


## Libraries that build on Neurgo

* [neurvolve](https://github.com/tleyden/neurvolve) builds on this library to support evolution-based learning.

## Related Work

[DXNN2](https://github.com/CorticalComputer/DXNN2) - Pure Erlang TPEULN (Topology & Parameter Evolving Universal Learning Network).  


## Related Publications

[Handbook of Neuroevolution Through Erlang](http://www.amazon.com/Handbook-Neuroevolution-Through-Erlang-Gene/dp/1461444624) _by Gene Sher_.


