package dla

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
	Origin
)

type TerrainType int

const (
	Undefined TerrainType = iota
	Ocean
	Coast
	Lake
	Land
)

type Tile struct {
	Parent   *Tile
	Children []*Tile
	Occupied bool
	x        int
	y        int
	Dir      Direction
	EndDist  int
	Height   float64
	Type     TerrainType
	Water    float64
}

func NewTile(x int, y int) Tile {
	var t Tile
	t.Occupied = false
	t.x = x
	t.y = y
	t.Parent = nil
	t.Children = make([]*Tile, 0)
	t.EndDist = -1
	t.Height = 0.0
	t.Water = 0.0
	return t
}

func (t *Tile) SetParent(Parent *Tile) {
	t.Parent = Parent
	Parent.Children = append(Parent.Children, t)

}

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

// TODO rename to highest
func (t *TileQueue) PopHighest() *Tile {
	if len(t.Tiles) <= 0 {
		return nil
	}
	lowest := t.Tiles[0]
	lowest_idx := 0
	for idx, current := range t.Tiles {
		if current.Priority > lowest.Priority {
			lowest = current
			lowest_idx = idx
		}
	}
	t.Remove(lowest_idx)
	return lowest.Tile
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
