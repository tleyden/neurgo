package neurgo

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"io"
	"log"
)

func (cortex Cortex) RenderSVG(writer io.Writer) {

	width := 500
	height := 500
	canvas := svg.New(writer)

	canvas.Start(width, height)

	canvas.Rect(0, 0, width, height, canvas.RGB(255, 255, 255))

	x := 100
	xDelta := 100
	yDelta := 100
	radius := 25
	neuronFill := "fill:black"
	actuatorFill := "fill:red"
	sensorFill := "fill:blue"

	layerToNodeIdMap := cortex.NodeIdLayerMap()
	layerIndexes := layerToNodeIdMap.Keys()
	for _, layerIndex := range layerIndexes {

		y := 100
		nodeIds := layerToNodeIdMap[layerIndex]
		layerIndexStr := fmt.Sprintf("%v", layerIndex)

		canvas.Text(x, y, layerIndexStr, "font-size:12;fill:black")
		y += yDelta

		for _, nodeId := range nodeIds {
			log.Printf("nodeId: %v", nodeId)
			log.Printf("x: %v, y: %v", x, y)

			switch nodeId.NodeType {
			case NEURON:
				canvas.Circle(x, y, radius, neuronFill)
			case ACTUATOR:
				canvas.Circle(x, y, radius, actuatorFill)
			case SENSOR:
				canvas.Circle(x, y, radius, sensorFill)
			}

			y += yDelta
		}

		x += xDelta
	}

	canvas.End()

}
