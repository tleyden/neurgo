package neurgo

import (
	"fmt"
)

type Any interface{}

func panicIfZero(x int) {
	if x == 0 {
		msg := fmt.Sprintf("not expecting this to be zero")
		panic(msg)
	}

}

func panicIfNil(any Any) {
	if any == nil {
		msg := fmt.Sprintf("not expecting this to be nil")
		panic(msg)

	}
}
