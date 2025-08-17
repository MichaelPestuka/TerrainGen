package dla

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
)

type Grid struct {
	width  int
	height int
	Tiles  [][]Tile
}

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
			if useRandom {
				if rand.Float64() > 0.7 {
					tile_slice[j].Occupied = true
				}
			}
		}
		g.Tiles[i] = tile_slice
	}
	return g
}

func (g Grid) Tile(x int, y int) *Tile {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return nil
	}
	return &(g.Tiles[x][y])
}

func (g Grid) SetOccupation(x int, y int, Occupation bool) {
	g.Tiles[x][y].Occupied = Occupation
}

func (g Grid) Occupation(x int, y int) bool {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return true
	}
	return g.Tiles[x][y].Occupied
}

func (g Grid) PrintDirections() {
	for y := range g.width {
		for x := range g.height {
			if g.Tile(x, y).Occupied {
				switch g.Tile(x, y).Dir {
				case Up:
					fmt.Printf("A")
				case Down:
					fmt.Printf("V")
				case Left:
					fmt.Printf("<")
				case Right:
					fmt.Printf(">")
				case Origin:
					fmt.Printf("O")
				default:
					fmt.Printf("#")
				}
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

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

func (g Grid) DiagonalNeighbors(t *Tile) []*Tile { // idk just includes diagonals
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
func (g Grid) MooreNeighbors(t *Tile) []*Tile { // idk just includes diagonals
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

func (g *Grid) UpscaleBy3() {

	upscaled := NewGrid(g.width*3, g.height*3, false)

	for x := range g.width {
		for y := range g.height {
			if !g.Tile(x, y).Occupied {
				continue
			}
			upscaled.Tiles[x*3+1][y*3+1].Occupied = true
			upscaled.Tiles[x*3+1][y*3+1].Dir = g.Tile(x, y).Dir

			switch g.Tile(x, y).Dir {
			case Up:
				upscaled.Tiles[x*3+1][y*3].Dir = Up
				upscaled.Tiles[x*3+1][y*3].Occupied = true
				upscaled.Tiles[x*3+1][y*3-1].Dir = Up
				upscaled.Tiles[x*3+1][y*3-1].Occupied = true

			case Down:
				upscaled.Tiles[x*3+1][y*3+2].Dir = Down
				upscaled.Tiles[x*3+1][y*3+2].Occupied = true
				upscaled.Tiles[x*3+1][y*3+3].Dir = Down
				upscaled.Tiles[x*3+1][y*3+3].Occupied = true

			case Left:
				upscaled.Tiles[x*3][y*3+1].Dir = Left
				upscaled.Tiles[x*3][y*3+1].Occupied = true
				upscaled.Tiles[x*3-1][y*3+1].Dir = Left
				upscaled.Tiles[x*3-1][y*3+1].Occupied = true

			case Right:
				upscaled.Tiles[x*3+2][y*3+1].Dir = Right
				upscaled.Tiles[x*3+2][y*3+1].Occupied = true
				upscaled.Tiles[x*3+3][y*3+1].Dir = Right
				upscaled.Tiles[x*3+3][y*3+1].Occupied = true
			}

			// if x-1 >= 0 && g.Tile(x-1, y).Occupied {
			// 	upscaled.SetOccupation(x*3, y*3+1, true)
			// }
			// if x+1 < g.width && g.Tile(x+1, y).Occupied {
			// 	upscaled.SetOccupation(x*3+2, y*3+1, true)
			// }
			// if y-1 >= 0 && g.Tile(x, y-1).Occupied {
			// 	upscaled.SetOccupation(x*3+1, y*3, true)
			// }
			// if y+1 < g.height && g.Tile(x, y+1).Occupied {
			// 	upscaled.SetOccupation(x*3+1, y*3+2, true)
			// }
		}
	}
	*g = upscaled
}

func (g Grid) ToFloatGrid(useEndDistance bool) FloatGrid {
	f := NewFloatGrid(g.width, g.height)
	for x := range g.width {
		for y := range g.height {
			if g.Tile(x, y).Occupied {
				if useEndDistance {
					if g.Tile(x, y).EndDist >= 0 {
						f.SetValue(x, y, 1.0-1.0/(1.0+float64(g.Tile(x, y).EndDist)))
					} else {
						f.SetValue(x, y, 0.0)
					}
				} else {
					f.SetValue(x, y, 1.0)
				}
			}
		}
	}
	return f
}

func (g *Grid) RunDLACycles(cycles int, max_steps int) {
	g.Tiles[g.width/2][g.height/2].Occupied = true
	g.Tiles[g.width/2][g.height/2].Dir = Origin
	for range cycles {
		// side := rand.Int() % 4
		var x, y int
		// switch side {
		// case 0:
		// 	x = 0
		// 	y = rand.Int() % g.height
		// case 1:
		// 	x = g.width - 1
		// 	y = rand.Int() % g.height
		// case 2:
		// 	x = rand.Int() % g.width
		// 	y = 0
		// case 3:
		// 	x = rand.Int() % g.width
		// 	y = g.height - 1
		// }

		x = rand.Int() % g.width
		y = rand.Int() % g.height
		for g.Tile(x, y).Occupied {
			x = rand.Int() % g.width
			y = rand.Int() % g.height
		}
		for range max_steps {
			dir := rand.Int() % 4
			dx := x
			dy := y
			switch dir {
			case 0:
				dx = min(g.width-1, dx+1)
			case 1:
				dx = max(0, dx-1)
			case 2:
				dy = min(g.height-1, dy+1)
			case 3:
				dy = max(0, dy-1)
			}
			if g.Tile(dx, dy).Occupied {
				if dx == x-1 {
					g.Tile(x, y).Dir = Left
				} else if dx == x+1 {
					g.Tile(x, y).Dir = Right
				} else if dy == y-1 {
					g.Tile(x, y).Dir = Up
				} else if dy == y+1 {
					g.Tile(x, y).Dir = Down
				}
				g.SetOccupation(x, y, true)
				g.Tile(x, y).SetParent(g.Tile(dx, dy))
				break
			}
			x = dx
			y = dy
		}
	}
}

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
				// img.Set(x, y, color.RGBA{R: 0, G: 0, B: uint8(128.0 + g.Tile(x, y).Height*128.0), A: 255})
			case Lake:
				img.Set(x, y, color.RGBA{R: 100, G: 100, B: 255, A: 255})
			case Shallows:
				img.Set(x, y, color.RGBA{R: 0, G: 255, B: 255, A: 255})
			}
		}
	}
	// DrawDepressions(img, g.FindDepressions())
	img.Set(0, 0, color.RGBA{255, 255, 255, 255})
	img.Set(g.width-1, g.height-1, color.RGBA{255, 255, 255, 255})
	return img
}

func (g *Grid) FindDepressions() []*Tile {
	depressions := make([]*Tile, 0)
	for x := range g.width {
		for y := range g.height {
			current := g.Tile(x, y)
			neighbors := g.Neighbors(current)
			isLower := true
			for _, n := range neighbors {
				if current.Height >= n.Height {
					isLower = false
				}
			}
			if isLower {
				depressions = append(depressions, current)
			}
		}
	}
	return depressions
}

func (g *Grid) FillDepressions(rise float64) {
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
			n.Height = math.Max(current.Height+rise, n.Height)
			closed[n.x][n.y] = true
			queue.Push(n, n.Height)
		}
		for _, n := range g.DiagonalNeighbors(current) {
			if closed[n.x][n.y] {
				continue
			}
			if current.Height >= 0.5 {
				n.Height = math.Max(current.Height+rise*1.414, n.Height)
			}
			closed[n.x][n.y] = true
			queue.Push(n, n.Height)
		}
	}
	// fmt.Printf("Cycles: %d\n", i)

}

func DrawDepressions(img *image.RGBA, depressions []*Tile) {
	for _, d := range depressions {
		img.Set(d.x, d.y, color.RGBA{255, 0, 255, 255})
	}
}

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
			// current.Height = oceanHeight // Flatten ocean
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

func (g *Grid) FindShallows(shallowsLength int) {
	for range shallowsLength {
		newGrid := NewGrid(g.width, g.height, false)
		for y := range g.height {
			copy(newGrid.Tiles[y], g.Tiles[y])
		}
		for x := range g.width {
			for y := range g.height {
				if g.Tile(x, y).Type == Ocean {
					for _, n := range g.Neighbors(g.Tile(x, y)) {
						if n.Type == Coast || n.Type == Shallows {
							newGrid.Tile(x, y).Type = Shallows
							break
						}
					}
				}
			}
		}
		g.Tiles = newGrid.Tiles
	}
}

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
