package dla

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
	Origin
)

type Tile struct {
	Parent   *Tile
	Children []*Tile
	Occupied bool
	x        int
	y        int
	Dir      Direction
	EndDist  int
}

func NewTile(x int, y int) Tile {
	var t Tile
	t.Occupied = false
	t.x = x
	t.y = y
	t.Children = make([]*Tile, 0)
	t.EndDist = -1
	return t
}

func (t *Tile) SetParent(Parent *Tile) {
	t.Parent = Parent
	Parent.Children = append(Parent.Children, t)

}
