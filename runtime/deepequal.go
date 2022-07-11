// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/assert"
)

//-------------------------------------------------------------------

type inProgressStack []pair

type pair struct {
	x     Value
	y     Value
	xTodo []Value
	yTodo []Value
}

const inProgressMax = 16

// for deepEqual xTodo and yTodo will always be the same length
// but for SuObject compare that isn't true

func (ip *inProgressStack) push(x, y Value, xTodo, yTodo []Value) {
	if len(*ip) >= inProgressMax {
		panic("object nesting overflow")
	}
	if len(xTodo) > 0 || len(yTodo) > 0 {
		*ip = append(*ip, pair{x: x, y: y, xTodo: xTodo, yTodo: yTodo})
	}
}

func (ip *inProgressStack) next() (x Value, y Value) {
	var p *pair
	for {
		if len(*ip) == 0 {
			return nil, nil
		}
		p = &(*ip)[len(*ip)-1]
		if len(p.xTodo) > 0 || len(p.yTodo) > 0 {
			break
		}
		*ip = (*ip)[:len(*ip)-1]
	}
	if len(p.xTodo) == 0 {
		x = nil
	} else {
		x = p.xTodo[0]
		p.xTodo = p.xTodo[1:]
	}
	if len(p.yTodo) == 0 {
		y = nil
	} else {
		y = p.yTodo[0]
		p.yTodo = p.yTodo[1:]
	}
	return
}

func (ip *inProgressStack) has(x, y Value) bool {
	for _, p := range *ip {
		if x == p.x && y == p.y {
			return true
		}
	}
	return false
}

//-------------------------------------------------------------------

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
	// starting these as nil greatly reduces allocation
	var inProgress inProgressStack
	var stack []Value
	for {
		if y == nil { // y didn't have a named member that x did
			return false
		}
		if x == y || inProgress.has(x, y) {
			goto endOfLoop
		}
		tx = order[x.Type()]
		ty = order[y.Type()]
		if tx != ty {
			return false
		}
		switch tx {
		case types.Object:
			xTodo, listSize, namedSize := xdeObject(
				ToContainer(x).ToObject(), &stack)
			yTodo, sizesEqual := ydeObject(
				ToContainer(y).ToObject(), &stack, listSize, namedSize)
			if !sizesEqual {
				return false
			}
			assert.That(len(xTodo) == len(yTodo))
			inProgress.push(x, y, xTodo, yTodo)
		case types.Instance:
			xTodo, size := xdeInstance(x.(*SuInstance), &stack)
			yTodo, sizesEqual := ydeInstance(y.(*SuInstance), &stack, size)
			if !sizesEqual {
				return false
			}
			assert.That(len(xTodo) == len(yTodo))
			inProgress.push(x, y, xTodo, yTodo)
		default:
			if !x.Equal(y) {
				return false
			}
		}
	endOfLoop:
		x, y = inProgress.next()
		if x == nil {
			return true
		}
	}
}

func xdeObject(x *SuObject, stack *[]Value) (todo []Value, listSize, namedSize int) {
	if x.RLock() {
		defer x.RUnlock()
	} else if x.named.Size() == 0 {
		// not concurrent, just list, don't need to copy
		return x.list, len(x.list), 0
	}
	listSize, namedSize = len(x.list), x.named.Size()
	start := len(*stack)
	xSize := namedSize + listSize
	expand(stack, namedSize+listSize+namedSize)
	i := start
	iter := x.named.Iter()
	for k, v := iter(); k != nil; k, v = iter() {
		(*stack)[i] = v.(Value)
		// temporarily store the member
		// where ydeObject will replace it with its value
		(*stack)[i+xSize] = k.(Value)
		i++
	}
	copy((*stack)[start+namedSize:], x.list)
	return (*stack)[start : start+xSize], listSize, namedSize
}

// expand increases the length of the stack by extra
// when possible it just reslices the stack
func expand(stack *[]Value, extra int) {
	desiredLen := len(*stack) + extra
	for desiredLen > cap(*stack) {
		// reslice to extend length to the end of capacity
		*stack = (*stack)[:cap(*stack)]
		// append should now grow the slice
		*stack = append(*stack, nil)
	}
	*stack = (*stack)[:desiredLen]
}

func ydeObject(y *SuObject, stack *[]Value, listSize, namedSize int) (todo []Value, sizesEqual bool) {
	concurrent := y.RLock()
	if concurrent {
		defer y.RUnlock()
	}
	if len(y.list) != listSize || y.named.Size() != namedSize {
		return nil, false
	}
	if !concurrent && namedSize == 0 {
		// don't need to copy
		return y.list, true
	}
	n := len(*stack)
	start := n - namedSize
	for i := 0; i < namedSize; i++ {
		p := &(*stack)[start+i]
		*p = y.namedGet(*p)
	}
	*stack = append(*stack, y.list...)
	return (*stack)[start:], true
}

func xdeInstance(x *SuInstance, stack *[]Value) (todo []Value, size int) {
	if x.Lock() {
		defer x.Unlock()
	}
	size = x.size()
	start := len(*stack)
	expand(stack, 2*size)
	i := start
	for k, v := range x.Data {
		(*stack)[i] = v.(Value)
		// temporarily store the member
		// where ydeInstance will replace it with its value
		(*stack)[i+size] = SuStr(k)
		i++
	}
	return (*stack)[start : start+size], size
}

func ydeInstance(y *SuInstance, stack *[]Value, size int) (todo []Value, sizesEqual bool) {
	if y.Lock() {
		defer y.Unlock()
	}
	if y.size() != size {
		return nil, false
	}
	n := len(*stack)
	start := n - size
	for i := 0; i < size; i++ {
		p := &(*stack)[start+i]
		*p = y.Data[string((*p).(SuStr))]
	}
	return (*stack)[start:], true
}
