package dla

import (
	"fmt"
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

func (g Grid) PrintGrid() {
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
