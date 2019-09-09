package runtime

// iList is a list of values
type iList struct {
	list []interface{}
}

type equable interface{ Equal(other interface{}) bool }

func (il *iList) Push(v interface{}) {
	il.list = append(il.list, v)
}

func (il *iList) Pop() {
	il.list = il.list[:len(il.list)-1]
}
func (il *iList) Has(v interface{}) bool {
	for _, x := range il.list {
		if v.(equable).Equal(x) {
			return true
		}
	}
	return false
}

func (il *iList) Remove(v interface{}) bool {
	for i, x := range il.list {
		if v.(equable).Equal(x) {
			copy(il.list[i:], il.list[i+1:])
			il.list[len(il.list)-1] = nil // for gc
			il.list = il.list[:len(il.list)-1]
			return true
		}
	}
	return false
}
