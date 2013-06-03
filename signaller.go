package neurgo

type Signaller interface {

	propagateSignal()

}

func Run(signaller Signaller) {

	signaller.propagateSignal()

}
