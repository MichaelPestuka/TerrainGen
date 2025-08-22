package dla

import "math"

type TerrainType int

const (
	Undefined TerrainType = iota
	Ocean
	Coast
	Lake
	Land
	Shallows
)

type Tile struct {
	x             int
	y             int
	Height        float64
	Type          TerrainType
	ShoreDistance float64
}

// New tile constructor
func NewTile(x int, y int) Tile {
	var t Tile
	t.x = x
	t.y = y
	t.Height = 0.0
	t.ShoreDistance = math.Inf(1)
	return t
}

// Tile priority queue based on height implementation
type TileQueueElement struct {
	Tile     *Tile
	Priority float64
}

type TileQueue struct {
	Tiles []*TileQueueElement
}

func (t *TileQueue) Remove(index int) {
	if len(t.Tiles) <= 0 || index >= len(t.Tiles) {
		return
	} else if len(t.Tiles) == 1 {
		t.Tiles = make([]*TileQueueElement, 0)
	} else {
		t.Tiles[index] = t.Tiles[len(t.Tiles)-1]
		t.Tiles = t.Tiles[:len(t.Tiles)-1]
	}

}

func (t *TileQueue) PopHighest() *Tile {
	if len(t.Tiles) <= 0 {
		return nil
	}
	highest := t.Tiles[0]
	highest_idx := 0
	for idx, current := range t.Tiles {
		if current.Priority > highest.Priority {
			highest = current
			highest_idx = idx
		}
	}
	t.Remove(highest_idx)
	return highest.Tile
}

func (t *TileQueue) PopLowest() *Tile {
	if len(t.Tiles) <= 0 {
		return nil
	}
	lowest := t.Tiles[0]
	lowest_idx := 0
	for idx, current := range t.Tiles {
		if current.Priority < lowest.Priority {
			lowest = current
			lowest_idx = idx
		}
	}
	t.Remove(lowest_idx)
	return lowest.Tile
}

func (t *TileQueue) Push(tile *Tile, priority float64) {
	t.Tiles = append(t.Tiles, &TileQueueElement{tile, priority})
}
