package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestRenderSVG(t *testing.T) {

	filename := "out.svg"
	outfile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := outfile.Close(); err != nil {
			panic(err)
		}
	}()

	xnorCortex := XnorCortex()
	xnorCortex.RenderSVG(outfile)

	content, err2 := ioutil.ReadFile(filename)
	if err != nil {
		panic(err2)
	}

	contentStr := string(content)
	assert.True(t, len(contentStr) > 0)

}
