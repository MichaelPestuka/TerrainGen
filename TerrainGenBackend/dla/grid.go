// Package for generating heightmaps

package dla

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/ojrac/opensimplex-go"
)

type Grid struct {
	width  int
	height int
	Tiles  [][]Tile
}

// Grid constructor
func NewGrid(width int, height int, useRandom bool) Grid {
	var g Grid
	g.width = width
	g.height = height
	g.Tiles = make([][]Tile, width)
	for i := range width {
		// g.Tiles[i] = make([]Tile, height)
		tile_slice := make([]Tile, height)
		for j := range height {
			tile_slice[j] = NewTile(i, j)
		}
		g.Tiles[i] = tile_slice
	}
	return g
}

// Returns pointer to tile at coords x, y
func (g Grid) Tile(x int, y int) *Tile {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return nil
	}
	return &(g.Tiles[x][y])
}

// Returns direct tile neighbors
func (g Grid) Neighbors(t *Tile) []*Tile {
	neighbors := make([]*Tile, 0)
	if g.Tile(t.x+1, t.y) != nil {
		neighbors = append(neighbors, g.Tile(t.x+1, t.y))
	}
	if g.Tile(t.x-1, t.y) != nil {
		neighbors = append(neighbors, g.Tile(t.x-1, t.y))
	}
	if g.Tile(t.x, t.y+1) != nil {
		neighbors = append(neighbors, g.Tile(t.x, t.y+1))
	}
	if g.Tile(t.x, t.y-1) != nil {
		neighbors = append(neighbors, g.Tile(t.x, t.y-1))
	}
	return neighbors
}

// Only diagonal neighbor tiles
func (g Grid) DiagonalNeighbors(t *Tile) []*Tile {
	neighbors := make([]*Tile, 0)
	if g.Tile(t.x+1, t.y+1) != nil {
		neighbors = append(neighbors, g.Tile(t.x+1, t.y+1))
	}
	if g.Tile(t.x-1, t.y+1) != nil {
		neighbors = append(neighbors, g.Tile(t.x-1, t.y+1))
	}
	if g.Tile(t.x+1, t.y-1) != nil {
		neighbors = append(neighbors, g.Tile(t.x+1, t.y-1))
	}
	if g.Tile(t.x-1, t.y-1) != nil {
		neighbors = append(neighbors, g.Tile(t.x-1, t.y-1))
	}
	return neighbors
}

// Returns direct and diagonal neighbors
func (g Grid) FullNeighbors(t *Tile) []*Tile {
	neighbors := make([]*Tile, 0)
	if g.Tile(t.x+1, t.y) != nil {
		neighbors = append(neighbors, g.Tile(t.x+1, t.y))
	}
	if g.Tile(t.x-1, t.y) != nil {
		neighbors = append(neighbors, g.Tile(t.x-1, t.y))
	}
	if g.Tile(t.x, t.y+1) != nil {
		neighbors = append(neighbors, g.Tile(t.x, t.y+1))
	}
	if g.Tile(t.x, t.y-1) != nil {
		neighbors = append(neighbors, g.Tile(t.x, t.y-1))
	}
	if g.Tile(t.x+1, t.y+1) != nil {
		neighbors = append(neighbors, g.Tile(t.x+1, t.y+1))
	}
	if g.Tile(t.x-1, t.y+1) != nil {
		neighbors = append(neighbors, g.Tile(t.x-1, t.y+1))
	}
	if g.Tile(t.x+1, t.y-1) != nil {
		neighbors = append(neighbors, g.Tile(t.x+1, t.y-1))
	}
	if g.Tile(t.x-1, t.y-1) != nil {
		neighbors = append(neighbors, g.Tile(t.x-1, t.y-1))
	}
	return neighbors
}

// Generates a terrain texture based on tile types
func (g Grid) TerrainTypeTexture() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
	for x := range g.width {
		for y := range g.height {
			switch g.Tile(x, y).Type {
			case Undefined:
				img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			case Land:
				img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			case Coast:
				img.Set(x, y, color.RGBA{R: 255, G: 255, B: 0, A: 255})
			case Ocean:
				img.Set(x, y, color.RGBA{R: 0, G: 0, B: 255, A: 255})
			case Lake:
				img.Set(x, y, color.RGBA{R: 100, G: 100, B: 255, A: 255})
			case Shallows:
				img.Set(x, y, color.RGBA{R: 0, G: 255, B: 255, A: 255})
			}
		}
	}
	return img
}

// Fills above water depresions
func (g *Grid) FillDepressions(rise float64) { // WARNING with randomness added, rivers cant be drawn inaccurately, remove before attempting
	closed := make([][]bool, g.width)
	for x := range g.width {
		closed[x] = make([]bool, g.height)
	}

	var queue TileQueue

	// Push edges
	for x := range g.width {
		queue.Push(g.Tile(x, 0), g.Tile(x, 0).Height)
		queue.Push(g.Tile(x, g.height-1), g.Tile(x, g.height-1).Height)
		closed[x][0] = true
		closed[x][g.height-1] = true
	}
	for y := range g.height {
		queue.Push(g.Tile(0, y), g.Tile(0, y).Height)
		queue.Push(g.Tile(g.width-1, y), g.Tile(g.width-1, y).Height)
		closed[0][y] = true
		closed[g.width-1][y] = true
	}
	i := 0
	for len(queue.Tiles) > 0 {
		i += 1
		current := queue.PopLowest()
		for _, n := range g.Neighbors(current) {
			if closed[n.x][n.y] {
				continue
			}
			n.Height = math.Max(current.Height+rise+rand.Float64()*(rise/10.0), n.Height) // RANDOMNESS ALERT
			closed[n.x][n.y] = true
			queue.Push(n, n.Height)
		}
		for _, n := range g.DiagonalNeighbors(current) {
			if closed[n.x][n.y] {
				continue
			}
			if current.Height >= 0.5 {
				n.Height = math.Max(current.Height+rise*1.414+rand.Float64()*(rise/10.0), n.Height) // RANDOMNESS ALERT
			}
			closed[n.x][n.y] = true
			queue.Push(n, n.Height)
		}
	}

}

// Sets tile types to ocean if connected to border with wiles below ocean height
func (g *Grid) DrawOcean(oceanHeight float64) {
	queue := make([]*Tile, 0)
	queue = append(queue, g.Tile(0, 0))
	queue = append(queue, g.Tile(g.width-1, g.height-1))
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Type != Undefined {
			continue
		}

		if current.Height <= oceanHeight {
			current.Type = Ocean
		} else {
			current.Type = Coast
			continue
		}

		if g.Tile(current.x+1, current.y) != nil && g.Tile(current.x+1, current.y).Type == Undefined {
			queue = append(queue, g.Tile(current.x+1, current.y))
		}
		if g.Tile(current.x-1, current.y) != nil && g.Tile(current.x-1, current.y).Type == Undefined {
			queue = append(queue, g.Tile(current.x-1, current.y))
		}
		if g.Tile(current.x, current.y+1) != nil && g.Tile(current.x, current.y+1).Type == Undefined {
			queue = append(queue, g.Tile(current.x, current.y+1))
		}
		if g.Tile(current.x, current.y-1) != nil && g.Tile(current.x, current.y-1).Type == Undefined {
			queue = append(queue, g.Tile(current.x, current.y-1))
		}
	}
}

// Changes underwater tile heights to a gentle slope, good for coastal wave rendering later
func (g *Grid) OceanSloping(step float64) {
	closed := make([][]bool, g.width)
	for x := range g.width {
		closed[x] = make([]bool, g.height)
	}
	open := make([]*Tile, 0)
	for x := range g.width {
		for y := range g.height {
			if g.Tile(x, y).Type == Coast {
				open = append(open, g.Tile(x, y))
				g.Tile(x, y).ShoreDistance = 0
			}
		}
	}
	for len(open) > 0 {
		closed[open[0].x][open[0].y] = true
		for _, n := range g.Neighbors(open[0]) {
			if !closed[n.x][n.y] {
				n.ShoreDistance = math.Min(n.ShoreDistance, open[0].ShoreDistance+1)
				closed[n.x][n.y] = true
				open = append(open, n)
			}
		}
		for _, n := range g.DiagonalNeighbors(open[0]) {
			if !closed[n.x][n.y] {
				n.ShoreDistance = math.Min(n.ShoreDistance, open[0].ShoreDistance+1.414)
				closed[n.x][n.y] = true
				open = append(open, n)
			}
		}
		if open[0].Type == Ocean {
			open[0].Height = 10 / (20 + open[0].ShoreDistance) // maxNeighbor
		}
		open = open[1:]
	}
}

// Applies a 2D bump function to make island roughly circular
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

// Fills grid with simplex noise
func (g *Grid) SimplexFill(octaves int, frequency float64) {
	noise := opensimplex.NewNormalized(rand.Int63())
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

// Normalizes tile heights to <0; 1>
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

// Exports tile heights to a float array for sending to client
func (g Grid) ExportHeights() []float64 {
	exp := make([]float64, g.width*g.height)
	for x := range g.width {
		for y := range g.height {
			exp[x+y*g.width] = g.Tile(x, y).Height
		}
	}
	return exp
}
