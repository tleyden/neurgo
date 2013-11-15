package neurgo

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"io"
	"os"
)

type Point struct {
	x int
	y int
}

type NodeCircleSVG struct {
	Point
	radius int
}

type NodeUUIDToCircleSVG map[string]NodeCircleSVG

func (cortex *Cortex) RenderSVGFile(filename string) {
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
}

func (cortex *Cortex) RenderSVG(writer io.Writer) {

	width := 1000
	height := 1000
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

			circleSVG := NodeCircleSVG{Point{x: x, y: y}, radius}
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

				layerDelta := tgtNodeId.LayerIndex - nodeId.LayerIndex
				if layerDelta > 0 {
					adjacent := layerToNodeIdMap.LayersAdjacent(nodeId.LayerIndex, tgtNodeId.LayerIndex)
					if adjacent {
						forwardConnectNodesSVG(canvas, srcCircle, tgtCircle)
					} else {
						forwardConnectDistantNodesSVG(canvas, srcCircle, tgtCircle)
					}

				} else if layerDelta == 0 {
					selfRecurrentConnectNodesSVG(canvas, srcCircle, tgtCircle)
				} else if layerDelta < 0 {
					recurrentConnectNodesSVG(canvas, srcCircle, tgtCircle)
				}

			}

		}

	}

}

func recurrentConnectNodesSVG(canvas *svg.SVG, src NodeCircleSVG, tgt NodeCircleSVG) {

	linestyle2 := []string{`stroke="turquoise"`, `stroke-linecap="round"`, `stroke-width="5"`, `fill="none"`}
	midpoint := midpoint(Point{x: src.x, y: src.y}, Point{x: tgt.x, y: tgt.y})
	controlX := midpoint.x
	controlY := midpoint.y - 50
	canvas.Qbez(src.x, src.y, controlX, controlY, tgt.x, tgt.y, linestyle2[0], linestyle2[1], linestyle2[2], linestyle2[3])

}

func selfRecurrentConnectNodesSVG(canvas *svg.SVG, src NodeCircleSVG, tgt NodeCircleSVG) {

	linestyle2 := []string{`stroke="turquoise"`, `stroke-linecap="round"`, `stroke-width="5"`, `fill="none"`}

	srcX := src.x - 10
	srcY := src.y
	tgtX := src.x + 10
	tgtY := src.y

	controlX := src.x
	controlY := src.y - 100
	canvas.Qbez(srcX, srcY, controlX, controlY, tgtX, tgtY, linestyle2[0], linestyle2[1], linestyle2[2], linestyle2[3])

}

func midpoint(p1 Point, p2 Point) Point {
	pResult := Point{}
	pResult.x = (p1.x + p2.x) / 2
	pResult.y = (p1.y + p2.y) / 2
	return pResult
}

func forwardConnectNodesSVG(canvas *svg.SVG, src NodeCircleSVG, tgt NodeCircleSVG) {
	linestyle := []string{`stroke="black"`, `stroke-linecap="round"`, `stroke-width="5"`}

	canvas.Line(src.x, src.y, tgt.x, tgt.y, linestyle[0], linestyle[1], linestyle[2])

}

func forwardConnectDistantNodesSVG(canvas *svg.SVG, src NodeCircleSVG, tgt NodeCircleSVG) {

	linestyle2 := []string{`stroke="black"`, `stroke-linecap="round"`, `stroke-width="5"`, `fill="none"`, `stroke-dasharray="10,10"`}
	midpoint := midpoint(Point{x: src.x, y: src.y}, Point{x: tgt.x, y: tgt.y})
	controlX := midpoint.x
	controlY := midpoint.y + 50
	canvas.Qbez(src.x, src.y, controlX, controlY, tgt.x, tgt.y, linestyle2[0], linestyle2[1], linestyle2[2], linestyle2[3], linestyle2[4])

}
