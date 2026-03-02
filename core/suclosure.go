// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"

	"github.com/apmckinlay/gsuneido/core/types"
)

// SuClosure is an instance of a closure block
type SuClosure struct {
	this Value
	// parent is the Frame of the outer function that created this closure.
	// It is used by interp to handle block returns.
	parent *Frame
	shared *Shared // captured shared variables from parent frame
	*SuFunc
}

// Value interface

var _ Value = (*SuClosure)(nil)

func (c *SuClosure) String() string {
	return strings.Replace(c.SuFunc.String(), "block */", "closure */", 1)
}

func (c *SuClosure) Equal(other any) bool {
	return c == other
}

func (c *SuClosure) Call(th *Thread, this Value, as *ArgSpec) Value {
	fn := c.SuFunc

	// normally done by SuFunc Call
	th.Args(&c.ParamSpec, as)
	return th.invokeClosure(fn, this, c)
}

func (*SuClosure) Type() types.Type {
	return types.Block
}

func (c *SuClosure) SetConcurrent() {
	if c.this != nil {
		c.this.SetConcurrent()
	}
	if c.shared == nil || c.shared.concurrent {
		return
	}
	c.shared.concurrent = true
	// make shared values concurrent
	for _, x := range c.shared.values {
		if x != nil {
			x.SetConcurrent()
		}
	}
}

func (c *SuClosure) IsConcurrent() Value {
	return SuBool(c.shared != nil && c.shared.concurrent)
}
