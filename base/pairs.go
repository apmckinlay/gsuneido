package base

// pairs is used by object and instance Equals
// to prevent infinite recursion from self references
type pairs []pair

type pair struct {
	x Value
	y Value
}

const maxpairs = 20

func newpairs() pairs {
	return make([]pair, 0, maxpairs)
}

func (ps *pairs) push(x Value, y Value) {
	if len(*ps) > maxpairs {
		panic("object equals nesting overflow")
	}
	*ps = append(*ps, pair{x, y})
}

func (ps pairs) contains(x Value, y Value) bool {
	for _, p := range ps {
		if p.x == x && p.y == y {
			return true
		}
	}
	return false
}
