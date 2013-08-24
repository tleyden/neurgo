package neurgo

import (
	"os"
	"testing"
)

func TestRenderSVG(t *testing.T) {

	outfile, err := os.Create("/Users/traun/tmp/out.svg")
	if err != nil {
		panic(err)
	}
	// close outfile on exit and check for its returned error
	defer func() {
		if err := outfile.Close(); err != nil {
			panic(err)
		}
	}()

	xnorCortex := XnorCortex()
	xnorCortex.RenderSVG(outfile)

}
