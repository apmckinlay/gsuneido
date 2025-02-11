// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
)

var _ = builtin(record, "(@args)")

func record(arg Value) Value {
	return SuRecordFromObject(arg.(*SuObject))
}

var _ = exportMethods(&RecordMethods, "record")

var _ = method(record_AttachRule, "(key,callable)")

func record_AttachRule(this, arg1, arg2 Value) Value {
	this.(*SuRecord).AttachRule(arg1, arg2)
	return nil
}

var _ = method(record_Clear, "()")

func record_Clear(this Value) Value {
	this.(*SuRecord).Clear()
	return nil
}

var _ = method(record_GetDeps, "(field)")

func record_GetDeps(this, arg Value) Value {
	return this.(*SuRecord).GetDeps(ToStr(arg))
}

var _ = method(record_Delete, "()")

func record_Delete(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	if as.Nargs != 0 {
		return ob_Delete(th, as, this, args)
	}
	trace.Dbms.Println("Record Drop", this)
	this.(*SuRecord).DbDelete(th)
	return nil
}

var _ = method(record_Drop, "()")

func record_Drop(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	trace.Dbms.Println("Record Drop", this)
	this.(*SuRecord).DbDelete(th)
	return nil
}

var _ = method(record_Invalidate, "(@args)")

func record_Invalidate(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	r := this.(*SuRecord)
	iter := NewArgsIter(as, args)
	for {
		k, v := iter()
		if k != nil || v == nil {
			break
		}
		r.Invalidate(th, AsStr(v))
	}
	return nil
}

var _ = method(record_NewQ, "()")

func record_NewQ(this Value) Value {
	return SuBool(this.(*SuRecord).IsNew())
}

var _ = method(record_Observer, "(observer)")

func record_Observer(this, arg Value) Value {
	this.(*SuRecord).Observer(arg)
	return nil
}

var _ = method(record_PreSet, "(field,value)")

func record_PreSet(this, arg1, arg2 Value) Value {
	this.(*SuRecord).PreSet(arg1, arg2)
	return nil
}

var _ = method(record_RemoveObserver, "(observer)")

func record_RemoveObserver(this, arg Value) Value {
	return SuBool(this.(*SuRecord).RemoveObserver(arg))
}

var _ = method(record_SetDeps, "(field,deps)")

func record_SetDeps(this, arg1, arg2 Value) Value {
	this.(*SuRecord).SetDeps(ToStr(arg1), ToStr(arg2))
	return nil
}

var _ = method(record_Transaction, "()")

func record_Transaction(this Value) Value {
	t := this.(*SuRecord).Transaction()
	if t == nil || t.Ended() {
		return False
	}
	return t
}

var _ = method(record_Update, "(record = false)")

func record_Update(th *Thread, this Value, args []Value) Value {
	trace.Dbms.Println("Record Update", this)
	this.(*SuRecord).DbUpdate(th, args[0])
	return nil
}
