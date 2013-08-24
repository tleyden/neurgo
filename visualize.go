package neurgo

import (
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

	xOffset := 50

	layerToNodeIdMap := cortex.NodeIdLayerMap()
	layerIndexes := layerToNodeIdMap.Keys()
	for _, layerIndex := range layerIndexes {

		yOffset := 50
		nodeIds := layerToNodeIdMap[layerIndex]
		for _, nodeId := range nodeIds {
			log.Printf("nodeId: %v", nodeId)
			canvas.Circle(xOffset, yOffset, 5, "fill:black")
			yOffset += 50
		}

		xOffset += 50
	}

	canvas.End()

}
