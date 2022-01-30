// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
)

var _ = builtin1("Record(@args)",
	func(arg Value) Value {
		return SuRecordFromObject(arg.(*SuObject))
	})

func init() {
	RecordMethods = Methods{
		"AttachRule": method2("(key,callable)", func(this, arg1, arg2 Value) Value {
			this.(*SuRecord).AttachRule(arg1, arg2)
			return nil
		}),
		"Clear": method0(func(this Value) Value {
			this.(*SuRecord).Clear()
			return nil
		}),
		"GetDeps": method1("(field)", func(this, arg Value) Value {
			return this.(*SuRecord).GetDeps(ToStr(arg))
		}),
		"Delete": methodRaw("()",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				k, v := NewArgsIter(as, args)()
				if k != nil || v != nil {
					return obDelete(t, as, this, args)
				}
				trace.Dbms.Println("Record Delete", this)
				this.(*SuRecord).DbDelete(t)
				return nil
			}),
		"Invalidate": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				r := this.(*SuRecord)
				iter := NewArgsIter(as, args)
				for {
					k, v := iter()
					if k != nil || v == nil {
						break
					}
					r.Invalidate(t, AsStr(v))
				}
				return nil
			}),
		"New?": method0(func(this Value) Value {
			return SuBool(this.(*SuRecord).IsNew())
		}),
		"Observer": method1("(observer)", func(this, arg Value) Value {
			this.(*SuRecord).Observer(arg)
			return nil
		}),
		"PreSet": method2("(field,value)", func(this, arg1, arg2 Value) Value {
			this.(*SuRecord).PreSet(arg1, arg2)
			return nil
		}),
		"RemoveObserver": method1("(observer)", func(this, arg Value) Value {
			return SuBool(this.(*SuRecord).RemoveObserver(arg))
		}),
		"SetDeps": method2("(field,deps)", func(this, arg1, arg2 Value) Value {
			this.(*SuRecord).SetDeps(ToStr(arg1), ToStr(arg2))
			return nil
		}),
		"Transaction": method0(func(this Value) Value {
			t := this.(*SuRecord).Transaction()
			if t == nil || t.Ended() {
				return False
			}
			return t
		}),
		"Update": method("(record = false)",
			func(t *Thread, this Value, args []Value) Value {
				trace.Dbms.Println("Record Update", this)
				this.(*SuRecord).DbUpdate(t, args[0])
				return nil
			}),
	}
}
