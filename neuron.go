package neurgo

type Neuron struct {
	Inbound  []InboundConnection
	Outbound []OutboundConnection
	Closing  chan bool
	Data     chan DataMessage
}
