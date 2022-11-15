package goevo

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"math/rand"
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
		NeuronSize:        15,
		InputNeuronColor:  color.RGBA{0, 255, 0, 255},
		HiddenNeuronColor: color.RGBA{255, 0, 255, 255},
		OutputNeuronColor: color.RGBA{255, 255, 0, 255},
	}
}

func (v *GenotypeVisualiser) DrawImage(g *Genotype) draw.Image {
	nodeYPosses := make(map[int]int)
	nodeXPosses := make(map[int]int)
	img := image.NewRGBA(image.Rect(0, 0, v.ImgSizeX, v.ImgSizeY))
	countsInp, countsHid, countsOut := g.Topology()
	for i := 0; i < countsInp; i++ {
		yPos := getPaddedPosition(i, countsInp, v.ImgSizeY, 0.8)
		nid := g.NeuronOrder[i]
		nodeXPosses[nid] = drawInputNeuron(img, v, yPos)
		nodeYPosses[nid] = yPos
	}
	for i := 0; i < countsOut; i++ {
		yPos := getPaddedPosition(i, countsOut, v.ImgSizeY, 0.8)
		nid := g.NeuronOrder[i+countsInp+countsHid]
		nodeXPosses[nid] = drawOutputNeuron(img, v, yPos)
		nodeYPosses[nid] = yPos
	}
	for i := 0; i < countsHid; i++ {
		posX := getPaddedPosition(i+1, countsHid+2, v.ImgSizeX, 0.8)
		avPos := 0
		avPosN := 0
		nid := g.NeuronOrder[i+countsInp]
		for cid := range g.Synapses {
			if g.Synapses[cid].To == nid {
				avPos += nodeYPosses[g.Synapses[cid].From]
				avPosN++
			}
		}
		yPos := avPos/avPosN + rand.Intn(100) - 50
		nodeYPosses[nid] = yPos
		nodeXPosses[nid] = posX
		drawHiddenNeuron(img, v, posX, yPos)
	}
	for cid := range g.Synapses {
		w := g.Synapses[cid].Weight
		isRecurrent := g.InverseNeuronOrder[g.Synapses[cid].From] > g.InverseNeuronOrder[g.Synapses[cid].To]
		drawConnection(img, nodeXPosses, nodeYPosses, g.Synapses[cid].From, g.Synapses[cid].To, w, v.NeuronSize, isRecurrent)
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
	for x := -r; x < r; x++ {
		for y := -r; y < r; y++ {
			if (x*x)+(y*y) <= r*r {
				img.Set(x0+x, y0+y, c)
			}
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

/*
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
*/
func thickLine(img draw.Image, x0, y0, x1, y1, w int, c color.Color) {
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
		st := int(w / 2)
		for xo := -st; xo < w-st; xo++ {
			for yo := -st; yo < w-st; yo++ {
				img.Set(x0+xo, y0+yo, c)
			}
		}
	}

}

func drawConnection(img draw.Image, xPoses, yPoses map[int]int, startID, endID int, w float64, r int, isRecurrent bool) {
	startX, startY := xPoses[startID]+r, yPoses[startID]
	endX, endY := xPoses[endID]-r, yPoses[endID]
	if w > 1 {
		w = 1
	} else if w < -1 {
		w = -1
	}
	//isRecurrent := startX > endX
	var col color.Color
	if w > 0 {
		if isRecurrent {
			col = color.RGBA{G: 255, A: 255}
		} else {
			col = color.RGBA{B: 255, A: 255}
		}
	} else {
		if isRecurrent {
			col = color.RGBA{G: 255, R: 255, A: 255}
		} else {
			col = color.RGBA{R: 255, A: 255}
		}
	}
	width := int(math.Max(math.Min(math.Abs(w), 1)*10, 1))
	thickLine(img, startX, startY, endX, endY, width, col)
}
