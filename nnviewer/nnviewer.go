package main

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
	"os"
)

func init() {
	logg.LogKeys["DEBUG"] = true
}

func main() {

	pathToCortex := os.Args[1]
	cortex := unmarshalCortex(pathToCortex)
	renderSVG(cortex)
	dumpCompactDescription(cortex)
	prettyPrintJSON(cortex)

}

func unmarshalCortex(pathToCortex string) *ng.Cortex {

	cortex, err := ng.NewCortexFromJSONFile(pathToCortex)
	if err != nil {
		logg.LogPanic("Error reading cortex from: %v.  Err: %v", pathToCortex, err)
	}
	return cortex

}

func dumpCompactDescription(cortex *ng.Cortex) {
	logg.LogTo("DEBUG", "Compact: %v", cortex.StringCompact())
}

func renderSVG(cortex *ng.Cortex) {

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

func prettyPrintJSON(cortex *ng.Cortex) {

	json, err := json.MarshalIndent(cortex, "", "    ")
	if err != nil {
		panic(err)
	}
	jsonString := fmt.Sprintf("%s", json)
	logg.LogTo("DEBUG", "%v", jsonString)

}
