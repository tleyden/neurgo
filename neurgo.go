
package neurgo

type VectorChannel chan []float64


type Connectable interface {  // TODO: move to connectable.go

	ConnectBidirectional(target Connectable)
	ConnectBidirectionalWeighted(target Connectable, weights []float64)

	connectOutboundWithChannel(target Connectable, channel VectorChannel) 
	connectInboundWithChannel(source Connectable, channel VectorChannel, weights []float64)

}

