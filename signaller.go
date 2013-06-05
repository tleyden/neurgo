package neurgo

type Signaller interface {

	propagateSignal()

}

func Run(signaller Signaller) {

	for {
		signaller.propagateSignal()
	}

}
