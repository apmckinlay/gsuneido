package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	RecordMethods = Methods{
		"Add": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				r := this.(*SuRecord)
				iter := NewArgsIter(as, args)
				if at := getNamed(as, args, SuStr("at")); at != nil {
					if i, ok := at.IfInt(); ok {
						addAt(ToObject(r), i, iter)
					} else {
						putAt(func(k Value, v Value) {
							r.Put(t, k, v)
						}, at, iter)
					}
				} else {
					addAt(ToObject(r), r.ListSize(), iter)
				}
				return this
			}),
		"AttachRule": method2("(key,callable)", func(this, arg1, arg2 Value) Value {
			this.(*SuRecord).AttachRule(arg1, arg2)
			return nil
		}),
		"Clear": method0(func(this Value) Value {
			this.(*SuRecord).Clear()
			return nil
		}),
		"Copy": method0(func(this Value) Value {
			return this.(*SuRecord).Copy()
		}),
		"Delete": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				r := this.(*SuRecord)
				if all := getNamed(as, args, SuStr("all")); all == True {
					r.Clear()
				} else {
					iter := NewArgsIter(as, args)
					for {
						k, v := iter()
						if k != nil || v == nil {
							break
						}
						r.Delete(t, v)
					}
				}
				return this
			}),
		"Erase": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				r := this.(*SuRecord)
				iter := NewArgsIter(as, args)
				for {
					k, v := iter()
					if k != nil || v == nil {
						break
					}
					r.Erase(t, v)
				}
				return this
			}),
		"GetDeps": method1("(field)", func(this, arg Value) Value {
			return this.(*SuRecord).GetDeps(IfStr(arg))
		}),
		"Invalidate": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				iter := NewArgsIter(as, args)
				for {
					k, v := iter()
					if k != nil || v == nil {
						break
					}
					this.(*SuRecord).Invalidate(ToStr(v))
				}

				return nil
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
			this.(*SuRecord).SetDeps(IfStr(arg1), IfStr(arg2))
			return nil
		}),
	}
}
