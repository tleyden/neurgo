package neurgo

type NodeId struct {
	Uuid     string
	NodeType string
}

type DataMessage struct {
	SenderId NodeId
	Inputs   []float32
}
