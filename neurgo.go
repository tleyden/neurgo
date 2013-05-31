
package neurgo

type VectorChannel chan []float32


type Connectable interface {

	ConnectBidirectionalWeighted(target Connectable, weights []float32)
	ConnectBidirectionalUnweighted(target Connectable)

	connectOutbound(target Connectable, channel VectorChannel) 
	connectInbound(source Connectable, channel VectorChannel, weights []float32)

}

