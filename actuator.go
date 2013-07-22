package neurgo

// TODO: need a "function" which is called to actuate based on data

type Actuator struct {
	Inbound      []InboundConnection
	Closing      chan bool
	Data         chan DataMessage
	VectorLength uint
}
