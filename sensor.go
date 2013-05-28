
package neurgo

type syncFunction func() []float32

type Sensor struct {
	SyncFunction syncFunction
}

func (sensor *Sensor) Connect_with_weights(target NeuralNode, weights []float32) {

}

func (sensor *Sensor) Connect(target NeuralNode) {

}

func (sensor *Sensor) DoSomething() {

}
