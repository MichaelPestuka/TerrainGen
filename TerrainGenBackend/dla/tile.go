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
}

func NewTile(x int, y int) Tile {
	var t Tile
	t.Occupied = false
	t.x = x
	t.y = y
	return t
}

func (t Tile) SetParent(Parent *Tile) {
	t.Parent = Parent
	Parent.Children = append(Parent.Children, &t)

}
