package goevo

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

type GenotypeVisualiser struct {
	ImgSizeX          int
	ImgSizeY          int
	NeuronSize        int
	InputNeuronColor  color.Color
	HiddenNeuronColor color.Color
	OutputNeuronColor color.Color
}

func NewGenotypeVisualiser() GenotypeVisualiser {
	return GenotypeVisualiser{
		ImgSizeX:          800,
		ImgSizeY:          800,
		NeuronSize:        13,
		InputNeuronColor:  color.RGBA{0, 255, 0, 255},
		HiddenNeuronColor: color.RGBA{255, 0, 255, 255},
		OutputNeuronColor: color.RGBA{255, 255, 0, 255},
	}
}

func (v *GenotypeVisualiser) DrawImage(g *Genotype) draw.Image {
	nodeYPosses := make(map[NodeID]int)
	nodeXPosses := make(map[NodeID]int)
	img := image.NewRGBA(image.Rect(0, 0, v.ImgSizeX, v.ImgSizeY))
	countsInp, countsHid, countsOut := g.GetNodeTypeCounts()
	for i := 0; i < countsInp; i++ {
		yPos := getPaddedPosition(i, countsInp, v.ImgSizeY, 0.8)
		nodeYPosses[g.Layers[i].ID] = yPos
		nodeXPosses[g.Layers[i].ID] = drawInputNeuron(img, v, yPos)
	}
	for i := 0; i < countsOut; i++ {
		yPos := getPaddedPosition(i, countsOut, v.ImgSizeY, 0.8)
		nodeXPosses[g.Layers[i+countsInp+countsHid].ID] = drawOutputNeuron(img, v, yPos)
		nodeYPosses[g.Layers[i+countsHid+countsInp].ID] = yPos
	}
	for i := 0; i < countsHid; i++ {
		posX := getPaddedPosition(i+1, countsHid+2, v.ImgSizeX, 0.8)
		avPos := 0
		avPosN := 0
		for c := range g.Connections {
			if g.Connections[c].Out == g.Layers[i+countsInp].ID {
				avPos += nodeYPosses[g.Connections[c].In]
				avPosN++
			}
		}
		yPos := avPos/avPosN + randRange(-50, 50)
		nodeYPosses[g.Layers[i+countsInp].ID] = yPos
		nodeXPosses[g.Layers[i+countsInp].ID] = posX
		drawHiddenNeuron(img, v, posX, yPos)
	}
	for i := range g.Connections {
		if g.Connections[i].Enabled {
			drawConnection(img, nodeXPosses, nodeYPosses, g.Connections[i])
		}
	}
	return img
}

func getPaddedPosition(n, max, width int, pad float64) int {
	if max == 1 {
		return width / 2
	}
	w := int(float64(width) * pad)
	p := (width - w) / 2
	i := w / (max - 1)
	return i*n + p
}

func (v *GenotypeVisualiser) DrawImageToJPGFile(filename string, g *Genotype) {
	img := v.DrawImage(g)
	f, err := os.Create(filename)
	if err != nil {
		panic(any(err))
	}
	defer f.Close()
	if err = jpeg.Encode(f, img, nil); err != nil {
		panic(any(err))
	}
}
func (v *GenotypeVisualiser) DrawImageToPNGFile(filename string, g *Genotype) {
	img := v.DrawImage(g)
	f, err := os.Create(filename)
	if err != nil {
		panic(any(err))
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		panic(any(err))
	}
}

func drawCircle(img draw.Image, x0, y0, r int, c color.Color) {
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)

	for x > y {
		img.Set(x0+x, y0+y, c)
		img.Set(x0+y, y0+x, c)
		img.Set(x0-y, y0+x, c)
		img.Set(x0-x, y0+y, c)
		img.Set(x0-x, y0-y, c)
		img.Set(x0-y, y0-x, c)
		img.Set(x0+y, y0-x, c)
		img.Set(x0+x, y0-y, c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}

func drawNeuron(img draw.Image, g *GenotypeVisualiser, posX, posY int, c color.Color) {
	drawCircle(img, posX, posY, g.NeuronSize, c)
	drawCircle(img, posX, posY, g.NeuronSize-1, c)
}
func drawInputNeuron(img draw.Image, g *GenotypeVisualiser, posY int) int {
	xpos := g.NeuronSize + 10
	drawNeuron(img, g, xpos, posY, g.InputNeuronColor)
	return xpos
}
func drawOutputNeuron(img draw.Image, g *GenotypeVisualiser, posY int) int {
	xpos := g.ImgSizeX - (g.NeuronSize + 10)
	drawNeuron(img, g, xpos, posY, g.OutputNeuronColor)
	return xpos
}

func drawHiddenNeuron(img draw.Image, g *GenotypeVisualiser, posX, posY int) {
	drawNeuron(img, g, posX, posY, g.HiddenNeuronColor)
}

func line(img draw.Image, x0, y0, x1, y1 int, c color.Color) {
	var dx = math.Abs(float64(x1 - x0))
	var dy = math.Abs(float64(y1 - y0))
	var err = dx - dy
	var sx, sy = 1, 1

	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}

	img.Set(x0, y0, c)
	for x0 != x1 || y0 != y1 {
		var e2 = 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
		img.Set(x0, y0, c)
	}

}

func drawConnection(img draw.Image, xPoses, yPoses map[NodeID]int, con *ConnectionGene) {
	startID := con.In
	startX, startY := xPoses[startID], yPoses[startID]
	endID := con.Out
	endX, endY := xPoses[endID], yPoses[endID]
	w := con.Weight
	if w > 1 {
		w = 1
	} else if w < -1 {
		w = -1
	}
	c := uint8(255 * (w/2 + 0.5))
	ic := 255 - c
	line(img, startX, startY, endX, endY, color.RGBA{ic, 0, c, 255})
}
