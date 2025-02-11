// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/ftsearch"
)

type suFtsearch struct {
	staticClass[suFtsearch]
}

func init() {
	Global.Builtin("Ftsearch", &suFtsearch{})
}

func (*suFtsearch) String() string {
	return "Ftsearch /* builtin class */"
}

func (sfs *suFtsearch) Equal(other any) bool {
	return sfs == other
}

func (*suFtsearch) Lookup(_ *Thread, method string) Value {
	return ftsearchMethods[method]
}

var ftsearchMethods = methods("ftsearch")

var _ = staticMethod(ftsearch_Create, "()")

func ftsearch_Create() Value {
	return &suFtsBuilder{b: ftsearch.NewBuilder()}
}

var _ = staticMethod(ftsearch_Load, "(data)")

func ftsearch_Load(data Value) Value {
	return newSuFtsIndex(ftsearch.Unpack(ToStr(data)))
}

var _ = staticMethod(ftsearch_Members, "()")

func ftsearch_Members() Value {
	return ftsearch_members
}

var ftsearch_members = methodList(ftsearchMethods)

func newSuFtsIndex(idx *ftsearch.Index) *suFtsIndex {
	var si suFtsIndex
	si.idx.Store(idx)
	return &si
}

//-------------------------------------------------------------------

type suFtsBuilder struct {
	ValueBase[suFtsBuilder]
	b *ftsearch.Builder
}

func (sfb *suFtsBuilder) String() string {
	return sfb.b.String()
}

func (sfb *suFtsBuilder) Equal(other any) bool {
	return sfb == other
}

func (*suFtsBuilder) Lookup(_ *Thread, method string) Value {
	return ftsBuilderMethods[method]
}

var ftsBuilderMethods = methods("ftsBuilder")

var _ = method(ftsBuilder_Add, "(id, title, text)")

func ftsBuilder_Add(this, id, title, text Value) Value {
	b := this.(*suFtsBuilder).b
	b.Add(ToInt(id), ToStr(title), ToStr(text))
	return nil
}

var _ = method(ftsBuilder_Index, "()")

func ftsBuilder_Index(this Value) Value {
	b := this.(*suFtsBuilder).b
	this.(*suFtsBuilder).b = ftsearch.NewBuilder()
	return newSuFtsIndex(b.ToIndex())
}

var _ = method(ftsBuilder_Pack, "()")

func ftsBuilder_Pack(this Value) Value {
	b := this.(*suFtsBuilder).b
	return SuStr(b.Pack())
}

//-------------------------------------------------------------------

type suFtsIndex struct {
	ValueBase[suFtsIndex]
	idx atomic.Pointer[ftsearch.Index]
}

func (sfi *suFtsIndex) get() *ftsearch.Index {
	idx := sfi.idx.Load()
	if idx == nil {
		panic("can't use ftsIndex during Update")
	}
	return idx
}

func (sfi *suFtsIndex) String() string {
	return sfi.get().String()
}

func (sfi *suFtsIndex) Equal(other any) bool {
	return sfi == other
}

func (*suFtsIndex) SetConcurrent() {
	// protected by atomic
}

func (*suFtsIndex) Lookup(_ *Thread, method string) Value {
	return ftsIndexMethods[method]
}

var ftsIndexMethods = methods("ftsIndex")

var _ = method(ftsIndex_Search, "(query, scores = false)")

func ftsIndex_Search(this, query, scores Value) Value {
	scors := ToBool(scores)
	idx := this.(*suFtsIndex).get()
	docScores := idx.Search(ToStr(query))
	list := make([]Value, len(docScores))
	for i, ds := range docScores {
		if scors {
			ob := &SuObject{}
			ob.Set(SuStr("id"), SuInt(ds.DocId))
			ob.Set(SuStr("score"), SuDnum{Dnum: dnum.FromFloat(ds.Score)})
			list[i] = ob
		} else {
			list[i] = IntVal(ds.DocId)
		}
	}
	return NewSuObject(list)
}

var _ = method(ftsIndex_Update, "(id, oldTitle, oldText, newTitle, newText)")

func ftsIndex_Update(_ *Thread, this Value, args []Value) Value {
	sfi := this.(*suFtsIndex)
	idx := sfi.idx.Swap(nil)
	if idx == nil {
		panic("concurrent ftsIndex.Update not allowed")
	}
	idx.Update(ToInt(args[0]), ToStr(args[1]), ToStr(args[2]), ToStr(args[3]),
		ToStr(args[4]))
	assert.That(sfi.idx.CompareAndSwap(nil, idx))
	return nil
}

var _ = method(ftsIndex_Pack, "()")

func ftsIndex_Pack(this Value) Value {
	sfi := this.(*suFtsIndex)
	return SuStr(sfi.get().Pack())
}

var _ = method(ftsIndex_WordInfo, "(word)")

func ftsIndex_WordInfo(this, word Value) Value {
	sfi := this.(*suFtsIndex)
	return SuStr(sfi.get().WordInfo(ToStr(word)))
}
