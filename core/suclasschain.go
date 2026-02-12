// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"github.com/apmckinlay/gsuneido/core/types"
)

// SuClassChain captures a class inheritance chain
// so that changes don't affect existing instances or active static method calls
// @immutable
type SuClassChain struct {
	ValueBase[*SuClassChain]
	classes []*SuClass
}

func (cc *SuClassChain) Class() *SuClass {
	return cc.classes[0]
}

func (cc *SuClassChain) Parents() []*SuClass {
	return cc.classes
}

// Value interface --------------------------------------------------

var _ Value = (*SuClassChain)(nil)
var _ Named = (*SuClassChain)(nil)
var _ Findable = (*SuClassChain)(nil)

func (cc *SuClassChain) String() string {
	return cc.Class().String()
}

func (cc *SuClassChain) Show() string {
	return cc.Class().Show()
}

func (*SuClassChain) Type() types.Type {
	return types.Class
}

func (cc *SuClassChain) Get(th *Thread, m Value) Value {
	return cc.Class().get1(th, cc, m, cc.classes)
}

func (cc *SuClassChain) Lookup(th *Thread, method string) Value {
	return cc.Class().lookup(th, method, cc.classes)
}

func (cc *SuClassChain) Call(th *Thread, this Value, as *ArgSpec) Value {
	if this == nil {
		this = cc
	}
	if f := cc.Class().get2(th, "CallClass", cc.classes); f != nil {
		return f.Call(th, this, as)
	}
	return cc.Class().New(th, as)
}

func (cc *SuClassChain) Equal(other any) bool {
	switch o := other.(type) {
	case *SuClassChain:
		return cc.Class() == o.Class()
	case *SuClass:
		return cc.Class() == o
	default:
		return false
	}
}

func (*SuClassChain) SetConcurrent() {
	// immutable so ok
}

func (*SuClassChain) IsConcurrent() Value {
	return EmptyStr
}

func (cc *SuClassChain) GetName() string {
	return cc.Class().Name
}

func (cc *SuClassChain) Finder(_ *Thread, fn func(v Value, mb *MemBase) Value) Value {
	for _, c := range cc.classes {
		if x := fn(c, &c.MemBase); x != nil {
			return x
		}
	}
	return nil
}

func (cc *SuClassChain) StartCoverage(count bool) {
	cc.Class().StartCoverage(count)
}

func (cc *SuClassChain) StopCoverage() *SuObject {
	return cc.Class().StopCoverage()
}
