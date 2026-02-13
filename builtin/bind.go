// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
)

var _ = builtin(Bind, "(@args)")

func Bind(args Value) Value {
	ob := args.(*SuObject)
	if ob.ListSize() == 0 {
		panic("usage: Bind(func, ...)")
	}
	fn := ob.PopFirst()
	return &suBind{fn: fn, args: ob}
}

type suBind struct {
	ValueBase[*suBind]
	fn   Value
	args *SuObject
}

var _ Value = (*suBind)(nil)

func (b *suBind) Call(th *Thread, _ Value, as *ArgSpec) Value {
	if as.Nargs == 0 {
		return th.PushCall(b.fn, nil, &ArgSpecEach0, b.args) // fast path
	}
	args := th.Args(&ParamSpecAt, as)
	ob := args[0].(*SuObject)
	i := 0
	iter := b.args.ArgsIter()
	for k, v := iter(); v != nil; k, v = iter() {
		if k == nil {
			ob.Insert(i, v)
			i++
		} else if !ob.HasKey(k) {
			ob.Set(k, v)
		}
	}
	return th.PushCall(b.fn, nil, &ArgSpecEach0, ob)
}

func (b *suBind) Equal(other any) bool {
	if other, ok := other.(*suBind); ok {
		return b.fn.Equal(other.fn) && b.args.Equal(other.args)
	}
	return false
}

func (*suBind) Type() types.Type {
	return types.BuiltinFunction
}

func (*suBind) Params() string {
	return "()"
}

func (*suBind) Lookup(_ *Thread, method string) Value {
	if m, ok := ParamsMethods[method]; ok {
		return m
	}
	return nil
}

func (b *suBind) SetConcurrent() {
	b.fn.SetConcurrent()
	if b.args != nil {
		b.args.SetConcurrent()
	}
}
