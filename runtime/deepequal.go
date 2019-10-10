package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type pairs []pair

type pair struct {
	x Value
	y Value
}

const initpairs = 12  // handles almost all cases
const maxpairs = 1024 // must be less than inProgressStack type max

func newpairs() pairs {
	return make([]pair, 0, initpairs)
}

func (ps *pairs) push(x Value, y Value) {
	if len(*ps) > maxpairs {
		panic("object equals overflow")
	}
	*ps = append(*ps, pair{x, y})
}

func (ps *pairs) pop() {
	*ps = (*ps)[:len(*ps)-1]
}

func (ps pairs) top() (x Value, y Value) {
	p := ps[len(ps)-1]
	return p.x, p.y
}

func (ps pairs) topIndex() int {
	return len(ps) - 1
}

// used by Compare
func (ps pairs) contains(x Value, y Value) bool {
	for _, p := range ps {
		if x == p.x && y == p.y { // NOTE: == not Equals
			return true
		}
	}
	return false
}

//-------------------------------------------------------------------

type inProgressStack []uint16

func (ip *inProgressStack) push(i int) {
	*ip = append(*ip, uint16(i))
}

func (ip *inProgressStack) pop() {
	*ip = (*ip)[:len(*ip)-1]
}

func (ip inProgressStack) top() int {
	return int(ip[len(ip)-1])
}

//-------------------------------------------------------------------

const (
	pair_continue = iota
	pair_equal
	pair_notequal
)

var order [types.N]types.Type

func init() {
	for i := types.Boolean; i < types.N; i++ {
		order[i] = i
	}
	order[types.Record] = types.Object
	order[types.Except] = types.String
}

// deepEqual is used by SuObject (and SuRecord, SuSequence) and SuInstance.
// Default values and rules are not applied.
// It is structured to only lock one object at a time to avoid deadlocks.
// This means the result is undefined if there are concurrent modifications.
// The locking is to prevent data race errors or corruption.
func deepEqual(x, y Value) bool {
	var tx, ty types.Type
	inProgress := make(inProgressStack, 0, 8) // 8 handles almost all cases
	stack := newpairs()
	stack.push(x, y)
	for {
		x, y = stack.top()
		inProgress.push(stack.topIndex())
		if x == y {
			goto endOfLoop
		}
		if y == nil { // missing named
			return false
		}
		switch checkRecursive(x, y, stack) {
		case pair_equal:
			goto endOfLoop
		case pair_notequal:
			return false
		}
		tx = order[x.Type()]
		ty = order[y.Type()]
		if tx != ty {
			return false
		}
		switch tx {
		case types.Object:
			listSize, namedSize := xdeObject(ToContainer(x).ToObject(), &stack)
			if !ydeObject(ToContainer(y).ToObject(), stack, listSize, namedSize) {
				return false // sizes unequal
			}
		case types.Instance:
			size := xdeInstance(x.(*SuInstance), &stack)
			if !ydeInstance(y.(*SuInstance), stack, size) {
				return false // sizes unequal
			}
		default:
			if !x.Equal(y) {
				return false
			}
		}
	endOfLoop:
		for stack.topIndex() == inProgress.top() {
			stack.pop()
			inProgress.pop()
			if len(stack) == 0 {
				return true
			}
		}
	}
}

func xdeObject(x *SuObject, ps *pairs) (listSize, namedSize int) {
	if x.Lock() {
		defer x.Unlock()
	}
	for _, v := range x.list {
		ps.push(v, nil)
	}
	iter := x.named.Iter()
	for k, v := iter(); k != nil; k, v = iter() {
		ps.push(v.(Value), k.(Value))
	}
	return len(x.list), x.named.Size()
}

func ydeObject(y *SuObject, ps pairs, listSize, namedSize int) bool {
	if y.Lock() {
		defer y.Unlock()
	}
	if len(y.list) != listSize || y.named.Size() != namedSize {
		return false
	}
	n := len(ps)
	base := n - namedSize - listSize
	for i := 0; i < listSize; i++ {
		ps[base+i].y = y.list[i]
	}
	for i := n - namedSize; i < n; i++ {
		ps[i].y = y.namedGet(ps[i].y)
	}
	return true
}

func xdeInstance(x *SuInstance, ps *pairs) int {
	if x.Lock() {
		defer x.Unlock()
	}
	for k, v := range x.Data {
		ps.push(v, SuStr(k))
	}
	return x.size()
}

func ydeInstance(y *SuInstance, ps pairs, size int) bool {
	if y.Lock() {
		defer y.Unlock()
	}
	if y.size() != size {
		return false
	}
	n := len(ps)
	for i := n - size; i < n; i++ {
		ps[i].y = y.Data[string(ps[i].y.(SuStr))]
	}
	return true
}

func checkRecursive(x, y Value, ps pairs) int {
	n := len(ps) - 1 // -1 to skip top which is x,y
	for i := 0; i < n; i++ {
		if x == ps[i].x {
			if y == ps[i].y {
				return pair_equal // comparison of this pair already pending
			}
			return pair_notequal // structure different - looping mismatch
		}
	}
	return pair_continue // normal case, no recursion
}
