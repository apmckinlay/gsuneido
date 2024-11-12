// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

// SuneidoObject is the global Suneido object.
// It is a separate type to allow Suneido.Parse and Suneido.Compile methods.
type SuneidoObject struct {
	SuObject
}

// SuneidoObjectMethods is initialized by builtin/suneido.go
var SuneidoObjectMethods Methods

func (so *SuneidoObject) Lookup(th *Thread, method string) Value {
	if m := SuneidoObjectMethods[method]; m != nil {
		return m
	}
	return so.SuObject.Lookup(th, method)
}
