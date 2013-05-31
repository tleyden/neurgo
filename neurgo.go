
package neurgo

type VectorChannel chan []float32


type Connectable interface {

	ConnectBidirectional(target Connectable, weights []float32)
	connectOutbound(target Connectable, channel VectorChannel) 
	connectInboundWeighted(source Connectable, channel VectorChannel, weights []float32)

	// Connect(target Connectable)
}

