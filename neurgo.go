
package neurgo

type NeuralNode interface {
	DoSomething()
}

type Connectable interface {
	Connect_with_weights(target Connectable, weights []float32)
	Connect(target Connectable)
}
