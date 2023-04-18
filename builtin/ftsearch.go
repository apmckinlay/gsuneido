// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
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

var ftsearchMethods = methods()

func (*suFtsearch) Lookup(_ *Thread, method string) Callable {
	return ftsearchMethods[method]
}

var _ = staticMethod(ftsearch_Create, "()")

func ftsearch_Create() Value {
	return &suFtsBuilder{b: ftsearch.NewBuilder()}
}

var _ = staticMethod(ftsearch_Load, "(data)")

func ftsearch_Load(data Value) Value {
	return &suFtsIndex{idx: ftsearch.Unpack(ToStr(data))}
}

//-------------------------------------------------------------------

type suFtsBuilder struct {
	ValueBase[suFtsBuilder]
	b *ftsearch.Builder
}

func (fb *suFtsBuilder) String() string {
	return fb.b.String()
}

func (*suFtsBuilder) Lookup(_ *Thread, method string) Callable {
	return ftsBuilderMethods[method]
}

var ftsBuilderMethods = methods()

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
	return &suFtsIndex{idx: b.ToIndex()}
}

var _ = method(ftsBuilder_Pack, "()")

func ftsBuilder_Pack(this Value) Value {
	b := this.(*suFtsBuilder).b
	return SuStr(b.Pack())
}

//-------------------------------------------------------------------

type suFtsIndex struct {
	ValueBase[suFtsIndex]
	idx *ftsearch.Index
}

func (fi *suFtsIndex) String() string {
	return fi.idx.String()
}

func (*suFtsIndex) Lookup(_ *Thread, method string) Callable {
	return ftsIndexMethods[method]
}

func (*suFtsIndex) SetConcurrent() {
	// read-only so ok
}

var ftsIndexMethods = methods()

var _ = method(ftsIndex_Search, "(query, scores = false)")

func ftsIndex_Search(this, query, scores Value) Value {
	scors := ToBool(scores)
	idx := this.(*suFtsIndex).idx
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
	idx := this.(*suFtsIndex).idx
	idx.Update(ToInt(args[0]), ToStr(args[1]), ToStr(args[2]), ToStr(args[3]),
		ToStr(args[4]))
	return nil
}

var _ = method(ftsIndex_Pack, "()")

func ftsIndex_Pack(this Value) Value {
	idx := this.(*suFtsIndex).idx
	return SuStr(idx.Pack())
}

var _ = method(ftsIndex_WordInfo, "(word)")

func ftsIndex_WordInfo(this, word Value) Value {
	idx := this.(*suFtsIndex).idx
	return SuStr(idx.WordInfo(ToStr(word)))
}
