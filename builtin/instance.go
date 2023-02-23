// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = exportMethods(&InstanceMethods)

var _ = method(instance_BaseQ, "(class)")

func instance_BaseQ(t *Thread, this Value, args []Value) Value {
	instance := this.(*SuInstance)
	if instance.Base() == args[0] { // handle anonymous base class
		return True
	}
	// else search by name to handle unload
	c, ok := args[0].(*SuClass)
	if !ok {
		return False
	}
	name := c.Name
	return nilToFalse(instance.Finder(t,
		func(v Value, _ *MemBase) Value {
			if c, ok := v.(*SuClass); ok && name == c.Name {
				return True
			}
			return nil
		}))
}

var _ = method(instance_Copy, "()")

func instance_Copy(this Value) Value {
	return this.(*SuInstance).Copy()
}

var _ = method(instance_Delete, "(@args)")

func instance_Delete(t *Thread, as *ArgSpec, this Value, args []Value) Value {
	if all := getNamed(as, args, SuStr("all")); all == True {
		this.(*SuInstance).Clear()
	} else {
		iter := NewArgsIter(as, args)
		for {
			k, v := iter()
			if k != nil || v == nil {
				break
			}
			this.(*SuInstance).Delete(v)
		}
	}
	return this
}

var _ = method(instance_ReadonlyQ, "()")

func instance_ReadonlyQ(this Value) Value {
	return False
}
