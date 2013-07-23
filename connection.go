package neurgo

type InboundConnection struct {
	NodeId  *NodeId
	Weights []float64
}

type OutboundConnection struct {
	NodeId *NodeId
	Data   chan DataMessage
}

type weightedInput struct {
	senderNodeId *NodeId
	weights      []float64
	inputs       []float64
}

func createEmptyWeightedInputs(inbound []*InboundConnection) []*weightedInput {

	weightedInputs := make([]*weightedInput, len(inbound))
	for i, inboundConnection := range inbound {
		weightedInput := &weightedInput{
			senderNodeId: inboundConnection.NodeId,
			weights:      inboundConnection.Weights,
			inputs:       nil,
		}
		weightedInputs[i] = weightedInput
	}
	return weightedInputs

}
