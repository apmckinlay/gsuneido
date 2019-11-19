package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin("Record(@args)",
	func(_ *Thread, args []Value) Value {
		return newRecord(args)
	})

func newRecord(args []Value) *SuRecord {
	return SuRecordFromObject(args[0].(*SuObject))
}

func init() {
	RecordMethods = Methods{
		"AttachRule": method2("(key,callable)", func(this, arg1, arg2 Value) Value {
			this.(*SuRecord).AttachRule(arg1, arg2)
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
				this.(*SuRecord).DbDelete()
				return nil
			}),
		"Invalidate": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
				iter := NewArgsIter(as, args)
				for {
					k, v := iter()
					if k != nil || v == nil {
						break
					}
					this.(*SuRecord).Invalidate(t, AsStr(v))
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
				this.(*SuRecord).DbUpdate(t, args[0])
				return nil
			}),
	}
}
