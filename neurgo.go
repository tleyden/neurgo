
package neurgo

type VectorChannel chan []float32

type Connectable interface {
	Connect_with_weights(target Connectable, weights []float32)
	Connect(target Connectable)
}
