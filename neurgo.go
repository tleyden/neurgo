
package neurgo

type VectorChannel chan []float32


type Connectable interface {  // TODO: move to connectable.go

	ConnectBidirectional(target Connectable)
	ConnectBidirectionalWeighted(target Connectable, weights []float32)

	connectOutboundWithChannel(target Connectable, channel VectorChannel) 
	connectInboundWithChannel(source Connectable, channel VectorChannel, weights []float32)

}

