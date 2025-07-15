package dla

import (
	"fmt"
	"math/rand/v2"
)

func (g *Grid) RunCrystalGrowth(cycles int, branchChance float64, starts int) {

	growthQueue := make([]*Tile, 0)

	for range starts {
		x := rand.Int() % g.width
		y := rand.Int() % g.height
		g.Tile(x, y).Occupied = true
		g.Tile(x, y).Dir = Origin

		growthQueue = append(growthQueue, g.Tile(x, y))
		growthQueue = append(growthQueue, g.Tile(x, y))
		growthQueue = append(growthQueue, g.Tile(x, y))
		growthQueue = append(growthQueue, g.Tile(x, y))
	}

	for len(growthQueue) > 0 {
		// fmt.Printf("Q: %d\n", len(growthQueue))
		x := growthQueue[0].x
		y := growthQueue[0].y

		growthQueue = growthQueue[1:]

		if rand.Float64() < branchChance {
			growthQueue = append(growthQueue, g.Tile(x, y))
		}

		newDir := (g.Tile(x, y).Dir + Direction(rand.Int()%3)) % 4
		checked := 0
		for !g.CheckSurrounding(x, y, newDir) {
			newDir += 1
			checked += 1
			if checked == 4 {
				break
			}
		}
		if checked == 4 {
			continue
		}
		switch newDir {
		case Left:
			g.SetOccupation(x+1, y, true)
			growthQueue = append(growthQueue, g.Tile(x+1, y))
			g.Tile(x+1, y).Dir = Left
			g.Tile(x+1, y).SetParent(g.Tile(x, y))
		case Right:
			g.SetOccupation(x-1, y, true)
			growthQueue = append(growthQueue, g.Tile(x-1, y))
			g.Tile(x-1, y).Dir = Right
			g.Tile(x-1, y).SetParent(g.Tile(x, y))
		case Up:
			g.SetOccupation(x, y+1, true)
			growthQueue = append(growthQueue, g.Tile(x, y+1))
			g.Tile(x, y+1).Dir = Up
			g.Tile(x, y+1).SetParent(g.Tile(x, y))
		case Down:
			g.SetOccupation(x, y-1, true)
			growthQueue = append(growthQueue, g.Tile(x, y-1))
			g.Tile(x, y-1).Dir = Down
			g.Tile(x, y-1).SetParent(g.Tile(x, y))
		}

	}
}

func (g *Grid) CheckSurrounding(x int, y int, dir Direction) bool {
	switch dir {
	case Left:
		return !g.Occupation(x+1, y) && !g.Occupation(x+1, y+1) && !g.Occupation(x+1, y-1) && !g.Occupation(x+2, y) && !g.Occupation(x+2, y+1) && !g.Occupation(x+2, y-1)
	case Right:
		return !g.Occupation(x-1, y) && !g.Occupation(x-1, y+1) && !g.Occupation(x-1, y-1) && !g.Occupation(x-2, y) && !g.Occupation(x-2, y+1) && !g.Occupation(x-2, y-1)
	case Up:
		return !g.Occupation(x, y+1) && !g.Occupation(x+1, y+1) && !g.Occupation(x-1, y+1) && !g.Occupation(x, y+2) && !g.Occupation(x+1, y+2) && !g.Occupation(x-1, y+2)
	case Down:
		return !g.Occupation(x, y-1) && !g.Occupation(x+1, y-1) && !g.Occupation(x-1, y-1) && !g.Occupation(x, y-2) && !g.Occupation(x+1, y-2) && !g.Occupation(x-1, y-2)
	default:
		return false
	}
}

func (g *Grid) CalculateEndDistance() {
	origins := make([]*Tile, 0)
	for x := range g.width {
		for y := range g.height {
			if g.Tile(x, y).Dir == Origin {
				origins = append(origins, g.Tile(x, y))
			}
		}
	}
	for idx, origin := range origins {
		origin.RecursiveDistance()
		fmt.Printf("Origin %d dist: %d\n", idx, origin.EndDist)
	}
}

func (t *Tile) RecursiveDistance() int {
	if len(t.Children) == 0 {
		t.EndDist = 0
		return 0
	}
	maxDist := 0
	for _, child := range t.Children {
		maxDist = max(maxDist, child.RecursiveDistance())
	}
	t.EndDist = maxDist + 1
	return maxDist + 1
}
