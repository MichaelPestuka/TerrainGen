package dla

import (
	"fmt"

	"github.com/larspensjo/Go-simplex-noise/simplexnoise"
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

func (f FloatGrid) Value(x int, y int) float64 {
	return f.values[x][y]
}

func (f FloatGrid) SetValue(x int, y int, value float64) {
	f.values[x][y] = value
}

func (f FloatGrid) Print() {
	for y := range f.height {
		for x := range f.width {
			if f.Value(x, y) < 0.01 {
				fmt.Printf("  .  ")
			} else {
				fmt.Printf(" %.1f ", f.Value(x, y))
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (f *FloatGrid) BoxBlur(kernelSize int, keepPeaks bool) {
	blurred := NewFloatGrid(f.width, f.height)
	for x := range f.width {
		for y := range f.height {
			sum := 0.0
			samples := 0
			for dx := -kernelSize; dx <= kernelSize; dx++ {

				if x+dx < 0 || x+dx >= f.width {
					samples += kernelSize*2 + 1
					continue
				}

				for dy := -kernelSize; dy <= kernelSize; dy++ {

					if y+dy < 0 || y+dy >= f.height {
						samples += 1
						continue
					}
					sum += float64(f.Value(x+dx, y+dy))
					samples += 1
				}
			}
			blurred.SetValue(x, y, sum/float64(samples))
			if keepPeaks && f.Value(x, y) > 0.9 {
				blurred.SetValue(x, y, f.Value(x, y))
				// fmt.Printf("height: %.2f\n", f.Value(x, y))
			}
		}
	}
	*f = blurred
}

func (f FloatGrid) ExportHeights() []float64 {
	exp := make([]float64, f.width*f.height)
	for x := range f.width {
		for y := range f.height {
			exp[x+y*f.width] = f.Value(x, y)
		}
	}
	return exp
}

func (f *FloatGrid) SimplexFill(octaves int) {
	for x := range f.width {
		for y := range f.height {
			simplexnoise.Noise2(float64(x), float64(y))
		}
	}
}
