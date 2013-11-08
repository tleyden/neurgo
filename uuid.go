package neurgo

import (
	"fmt"
)

func NewUuid() string {
	// TODO: do a real uuid - https://github.com/nu7hatch/gouuid
	randInt := RandomIntInRange(0, 10000000000)
	return fmt.Sprintf("%d", randInt)
}
