package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin("Record(@args)",
	func(_ *Thread, args []Value) Value {
		return newRecord(args)
	})

func newRecord(args []Value) *SuRecord {
	ob := args[0].(*SuObject)
	ob.SetDefault(EmptyStr)
	return SuRecordFromObject(ob)
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
					this.(*SuRecord).Invalidate(AsStr(v))
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
			return this.(*SuRecord).Transaction()
		}),
		"Update": method("(record = false)",
			func(t *Thread, this Value, args []Value) Value {
				r := this.(*SuRecord)
				var ob Container = r
				if args[0] != False {
					ob = ToContainer(args[0])
				}
				r.DbUpdate(t, ob)
				return nil
			}),
	}
}
