
package neurgo

type VectorChannel chan []float32


type Connectable interface {

	ConnectBidirectional(target Connectable)
	ConnectBidirectionalWeighted(target Connectable, weights []float32)

	connectOutbound(target Connectable, channel VectorChannel) 
	connectInbound(source Connectable, channel VectorChannel, weights []float32)

}

