package neurgo

import (
	"fmt"
)

func NewUuid() string {
	// TODO: do a real uuid
	randInt := RandomIntInRange(0, 10000000000)
	return fmt.Sprintf("todo=%v", randInt)
}
