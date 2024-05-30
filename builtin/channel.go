// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
)

type suChannel struct {
	ValueBase[suChannel]
	concurrent bool
	ch         chan Value
}

var sc_timeout = 10 * time.Second

var _ = builtin(Channel, "(size = 4)")

func Channel(size Value) Value {
	return &suChannel{ch: make(chan Value, IfInt(size))}
}

var suChannelMethods = methods()

var _ = method(chan_Send, "(value)")

func chan_Send(this, val Value) Value {
	defer func() {
		if r := recover(); r != nil {
			// e.g. send on a closed channel
			panic(fmt.Sprint("Channel: ", r))
		}
	}()
	sc := this.(*suChannel)
	select {
	case sc.ch <- val:
		// value sent
	case <-time.After(sc_timeout):
		panic("Channel: Send timeout")
	}
	return nil
}

var _ = method(chan_Recv, "()")

func chan_Recv(this Value) Value {
	sc := this.(*suChannel)
	select {
	case val := <-sc.ch:
		if val == nil {
			return this // closed
		}
		if sc.concurrent {
			val.SetConcurrent()
		}
		return val
	case <-time.After(sc_timeout):
		panic("Channel: Recv timeout")
	}
}

var _ = method(chan_Recv2, "(channel)")

func chan_Recv2(this, arg Value) Value {
	sc := this.(*suChannel)
	sc2 := arg.(*suChannel)
	ob := &SuObject{}
	select {
	case val := <-sc.ch:
		ob.Add(Zero)
		if val != nil {
			if sc.concurrent {
				val.SetConcurrent()
			}
			ob.Add(val)
		}
	case val := <-sc2.ch:
		ob.Add(One)
		if val != nil {
			if sc2.concurrent {
				val.SetConcurrent()
			}
			ob.Add(val)
		}
	case <-time.After(sc_timeout):
		panic("Channel: Recv2 timeout")
	}
	return ob
}

var _ = method(chan_Close, "()")

func chan_Close(th *Thread, this Value, args []Value) Value {
	// WARNING there is a race if Send & Close are called concurrently
	defer func() {
		if r := recover(); r != nil {
			// e.g. close of a closed channel
			panic(fmt.Sprint("Channel: ", r))
		}
	}()
	sc := this.(*suChannel)
	close(sc.ch)
	return nil
}

// Value implementation

var _ Value = (*suChannel)(nil)

func (sc *suChannel) Equal(other any) bool {
	return sc == other
}

func (*suChannel) Lookup(_ *Thread, method string) Callable {
	return suChannelMethods[method]
}

func (sc *suChannel) SetConcurrent() {
	sc.concurrent = true
}
