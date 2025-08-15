package dla

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/ojrac/opensimplex-go"
)

type FloatGrid struct {
	width  int
	height int
	values [][]float64
}

func NewFloatGrid(width int, height int) FloatGrid {
	var f FloatGrid
	f.width = width
	f.height = height
	f.values = make([][]float64, width)
	for i := range width {
		// g.Tiles[i] = make([]Tile, height)
		value_slice := make([]float64, height)
		for j := range height {
			value_slice[j] = 0
		}
		f.values[i] = value_slice
	}
	return f
}

func (g FloatGrid) Value(x int, y int) float64 {
	return g.values[x][y]
}

func (g FloatGrid) SetValue(x int, y int, value float64) {
	g.values[x][y] = value
}

func (g Grid) PrintHeights() {
	for y := range g.height {
		for x := range g.width {
			if g.Tile(x, y).Height < 0.01 {
				fmt.Printf("  .  ")
			} else {
				fmt.Printf(" %.1f ", g.Tile(x, y).Height)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (g *FloatGrid) BoxBlur(kernelSize int, keepPeaks bool) {
	blurred := NewFloatGrid(g.width, g.height)
	for x := range g.width {
		for y := range g.height {
			sum := 0.0
			samples := 0
			for dx := -kernelSize; dx <= kernelSize; dx++ {

				if x+dx < 0 || x+dx >= g.width {
					samples += kernelSize*2 + 1
					continue
				}

				for dy := -kernelSize; dy <= kernelSize; dy++ {

					if y+dy < 0 || y+dy >= g.height {
						samples += 1
						continue
					}
					sum += float64(g.Value(x+dx, y+dy))
					samples += 1
				}
			}
			blurred.SetValue(x, y, sum/float64(samples))
			if keepPeaks && g.Value(x, y) > 0.9 {
				blurred.SetValue(x, y, g.Value(x, y))
				// fmt.Printf("height: %.2f\n", f.Value(x, y))
			}
		}
	}
	*g = blurred
}

func (g Grid) ExportHeights() []float64 {
	exp := make([]float64, g.width*g.height)
	for x := range g.width {
		for y := range g.height {
			exp[x+y*g.width] = g.Tile(x, y).Height
		}
	}
	return exp
}

func (g *Grid) Normalize() {
	max := math.Inf(-1)
	min := math.Inf(1)
	for x := range g.width {
		for y := range g.height {
			min = math.Min(g.Tile(x, y).Height, min)
			max = math.Max(g.Tile(x, y).Height, max)
		}
	}
	max -= min
	for x := range g.width {
		for y := range g.height {
			g.Tile(x, y).Height -= min
			g.Tile(x, y).Height /= max
		}
	}
}

func (g *Grid) CircleFilter(edgeOffsetPercent float64, slope float64) {

	smaller := g.width
	if g.width > g.height {
		smaller = g.height
	}
	edgeOffset := int(smaller * int(edgeOffsetPercent))
	for x := range g.width {
		for y := range g.height {
			squares := math.Pow(float64(x-g.width/2), 2) + math.Pow(float64(y-g.height/2), 2)
			if squares <= math.Pow(float64(g.width/2-edgeOffset), 2) {

				g.Tile(x, y).Height *= math.Exp(-(slope / (float64(smaller/2-edgeOffset) - math.Pow(squares, 0.5))))
			} else {
				g.Tile(x, y).Height *= 0.0
			}
		}
	}
}

func (g *Grid) SimplexFill(octaves int, frequency float64) {
	noise := opensimplex.NewNormalized(rand.Int64())
	amplitude := 1.0
	for range octaves {
		for x := range g.width {
			for y := range g.height {
				g.Tile(x, y).Height += amplitude * noise.Eval2(float64(x)*frequency/float64(g.width), float64(y)*frequency/float64(g.height))
			}
		}
		frequency *= 2.0
		amplitude /= 2.0
	}
}
