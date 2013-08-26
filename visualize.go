package neurgo

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"io"
	"log"
)

type NodeCircleSVG struct {
	x      int
	y      int
	radius int
}

type NodeUUIDToCircleSVG map[string]NodeCircleSVG

func (cortex *Cortex) RenderSVG(writer io.Writer) {

	width := 500
	height := 500
	x := 100
	xDelta := 100
	yDelta := 100
	radius := 25
	neuronFill := "fill:blue"
	actuatorFill := "fill:magenta"
	sensorFill := "fill:green"

	canvas := svg.New(writer)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, canvas.RGB(255, 255, 255))

	nodeUUIDToCircleSVG := make(NodeUUIDToCircleSVG)
	layerToNodeIdMap := cortex.NodeIdLayerMap()
	layerIndexes := layerToNodeIdMap.Keys()

	for _, layerIndex := range layerIndexes {

		y := 100
		nodeIds := layerToNodeIdMap[layerIndex]
		layerIndexStr := fmt.Sprintf("%v", layerIndex)

		canvas.Text(x, y, layerIndexStr, "font-size:12;fill:black")
		y += yDelta

		for _, nodeId := range nodeIds {

			switch nodeId.NodeType {
			case NEURON:
				canvas.Circle(x, y, radius, neuronFill)
			case ACTUATOR:
				canvas.Circle(x, y, radius, actuatorFill)
			case SENSOR:
				canvas.Circle(x, y, radius, sensorFill)
			}

			circleSVG := NodeCircleSVG{x: x, y: y, radius: radius}
			nodeUUIDToCircleSVG[nodeId.UUID] = circleSVG

			y += yDelta
		}

		x += xDelta
	}

	addConnectionsToSVG(cortex, canvas, nodeUUIDToCircleSVG)

	canvas.End()

}

func addConnectionsToSVG(cortex *Cortex, canvas *svg.SVG, nodeUUIDToCircleSVG NodeUUIDToCircleSVG) {

	layerToNodeIdMap := cortex.NodeIdLayerMap()
	layerIndexes := layerToNodeIdMap.Keys()

	// loop over all layers
	for _, layerIndex := range layerIndexes {

		nodeIds := layerToNodeIdMap[layerIndex]

		// loop over all nodes
		for _, nodeId := range nodeIds {
			log.Printf("nodeId: %v", nodeId)

			// lookup node (assuming it is an OutboundConnector)
			node := cortex.FindConnector(nodeId)
			if node == nil {
				// if not, ignore it (eg, actuator)
				continue
			}

			// loop over all outbound connections
			for _, outbound := range node.outbound() {
				tgtNodeId := outbound.NodeId
				srcCircle := nodeUUIDToCircleSVG[nodeId.UUID]
				tgtCircle := nodeUUIDToCircleSVG[tgtNodeId.UUID]
				forwardConnectNodesSVG(canvas, srcCircle, tgtCircle)
				recurrentConnectNodesSVG(canvas, srcCircle, tgtCircle)

			}

		}

	}

}

func recurrentConnectNodesSVG(canvas *svg.SVG, src NodeCircleSVG, tgt NodeCircleSVG) {

	// draw a bezier curve

	linestyle2 := []string{`stroke="turquoise"`, `stroke-linecap="round"`, `stroke-width="5"`, `fill="none"`}

	// line: src.x 100, src.y: 200 tgt.x: 200 tgt.y: 200
	srcX := src.x
	srcY := src.y
	tgtX := tgt.x
	tgtY := tgt.y
	controlX := (srcX + tgtX) / 2
	controlY := srcY - 50
	canvas.Qbez(srcX, srcY, controlX, controlY, tgtX, tgtY, linestyle2[0], linestyle2[1], linestyle2[2], linestyle2[3])

	// find the x midpoint
	// xMidpoint := (src.x + src.y) / 2

}

func forwardConnectNodesSVG(canvas *svg.SVG, src NodeCircleSVG, tgt NodeCircleSVG) {
	linestyle := []string{`stroke="black"`, `stroke-linecap="round"`, `stroke-width="5"`}

	log.Printf("line: src.x %v, src.y: %v tgt.x: %v tgt.y: %v", src.x, src.y, tgt.x, tgt.y)

	canvas.Line(src.x, src.y, tgt.x, tgt.y, linestyle[0], linestyle[1], linestyle[2])

}
