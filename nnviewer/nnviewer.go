package main

import (
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
	"os"
)

func init() {
	logg.LogKeys["DEBUG"] = true
}

func main() {

	pathToCortex := os.Args[1]
	logg.LogTo("DEBUG", "Opening file: %v", pathToCortex)

	cortex, err := ng.NewCortexFromJSONFile(pathToCortex)
	if err != nil {
		logg.LogPanic("Error reading cortex from: %v.  Err: %v", pathToCortex, err)
	}

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

	cortex.RenderSVG(outfile)
	logg.LogTo("DEBUG", "svg available here: %v", filename)

}
