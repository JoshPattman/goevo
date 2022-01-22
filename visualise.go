package goevo

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
)

type GenotypeVisualiser struct {
	ImgSizeX   int
	ImgSizeY   int
	NeuronSize int
}

func NewGenotypeVisualiser() GenotypeVisualiser {
	return GenotypeVisualiser{
		ImgSizeX:   800,
		ImgSizeY:   800,
		NeuronSize: 13,
	}
}

func (v *GenotypeVisualiser) DrawImage(g *Genotype) draw.Image {
	nodeYPosses := make(map[NodeID]int)
	nodeXPosses := make(map[NodeID]int)
	img := image.NewRGBA(image.Rect(0, 0, v.ImgSizeX, v.ImgSizeY))
	countsInp, countsHid, countsOut := g.GetNodeTypeCounts()
	for i := 0; i < countsInp; i++ {
		var pos float64
		if countsInp == 1 {
			pos = 0.5
		} else {
			pos = float64(i) / float64(countsInp-1)
		}
		w := float64(v.ImgSizeX) * 0.8
		pad := (v.ImgSizeX - int(w)) / 2
		yPos := int(pos*w) + pad
		nodeYPosses[g.Nodes[i].ID] = yPos
		nodeXPosses[g.Nodes[i].ID] = drawInputNeuron(img, v, yPos)
	}
	for i := 0; i < countsOut; i++ {
		var pos float64
		if countsOut == 1 {
			pos = 0.5
		} else {
			pos = float64(i) / float64(countsOut-1)
		}
		w := float64(v.ImgSizeX) * 0.8
		pad := (v.ImgSizeX - int(w)) / 2
		yPos := int(pos*w) + pad
		nodeXPosses[g.Nodes[i+countsInp+countsHid].ID] = drawOutputNeuron(img, v, yPos)
		nodeYPosses[g.Nodes[i+countsHid+countsInp].ID] = yPos
	}
	for i := 0; i < countsHid; i++ {
		pos := float64(i+1) / float64(countsHid+1)
		w := float64(v.ImgSizeX) * 0.8
		pad := (v.ImgSizeX - int(w)) / 2
		posX := int(pos*w) + pad
		avPos := 0
		avPosN := 0
		for c := range g.Connections {
			if g.Connections[c].Out == g.Nodes[i+countsInp].ID {
				avPos += nodeYPosses[g.Connections[c].In]
				avPosN++
			}
		}
		yPos := avPos/avPosN + randRange(-50, 50)
		nodeYPosses[g.Nodes[i+countsInp].ID] = yPos
		nodeXPosses[g.Nodes[i+countsInp].ID] = posX
		drawHiddenNeuron(img, v, posX, yPos)
	}
	for i := range g.Connections {
		if g.Connections[i].Enabled {
			startID := g.Connections[i].In
			startX, startY := nodeXPosses[startID], nodeYPosses[startID]
			endID := g.Connections[i].Out
			endX, endY := nodeXPosses[endID], nodeYPosses[endID]
			c := uint8(255 * (g.Connections[i].Weight/2 + 0.5))
			ic := 255 - c
			line(img, startX, startY, endX, endY, color.RGBA{ic, c / 2, c, 255})
		}
	}
	return img
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
}
func drawInputNeuron(img draw.Image, g *GenotypeVisualiser, posY int) int {
	xpos := g.NeuronSize + 10
	drawNeuron(img, g, xpos, posY, color.RGBA{0, 255, 0, 255})
	return xpos
}
func drawOutputNeuron(img draw.Image, g *GenotypeVisualiser, posY int) int {
	xpos := g.ImgSizeX - (g.NeuronSize + 10)
	drawNeuron(img, g, xpos, posY, color.RGBA{255, 255, 0, 255})
	return xpos
}

func drawHiddenNeuron(img draw.Image, g *GenotypeVisualiser, posX, posY int) {
	drawNeuron(img, g, posX, posY, color.RGBA{255, 0, 255, 255})
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
