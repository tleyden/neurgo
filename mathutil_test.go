package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
)

func TestSigmoid(t *testing.T) {

	assert.Equals(t, sigmoid(0), 0.5)

}
