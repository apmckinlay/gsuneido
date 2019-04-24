package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/pack"
)

//TODO rules

// SuRecord is an SuObject with observers and rules and a default value of ""
type SuRecord struct {
	SuObject
	observers       ValueList
	activeObservers ActiveList
}

//go:generate genny -in ../../GoTemplates/list/list.go -out alist.go -pkg runtime gen "V=active"
//go:generate genny -in ../../GoTemplates/list/list.go -out vlist.go -pkg runtime gen "V=Value"

func NewSuRecord() *SuRecord {
	return &SuRecord{SuObject: SuObject{defval: EmptyStr}}
}

func (r *SuRecord) Copy() *SuRecord {
	return &SuRecord{SuObject: *r.SuObject.Copy()}
}

func (*SuRecord) Type() types.Type {
	return types.Record
}

func (r *SuRecord) String() string {
	s := r.SuObject.String()
	return "[" + s[2:len(s)-1] + "]"
}

func (r *SuRecord) Show() string {
	s := r.SuObject.Show()
	return "[" + s[2:len(s)-1] + "]"
}

func (r *SuRecord) Put(t *Thread, key Value, val Value) {
	r.SuObject.Set(key, val)
	r.callObservers(t, key)
}

func (r *SuRecord) PreSet(key Value, val Value) {
	r.SuObject.Set(key, val)
}

func (r *SuRecord) Observer(ofn Value) {
	r.observers.Push(ofn)
}

func (r *SuRecord) RemoveObserver(ofn Value) bool {
	return r.observers.Remove(ofn)
}

func (r *SuRecord) callObservers(t *Thread, key Value) {
	for _, x := range r.observers.list {
		ofn := x.(Value)
		if !r.activeObservers.Has(active{ofn, key}) {
			func(ofn Value, key Value) {
				r.activeObservers.Push(active{ofn, key})
				defer r.activeObservers.Pop()
				t.CallAsMethod(r, ofn, key)
			}(ofn, key)
		}
	}
}

type active struct {
	obs Value
	key Value
}

func (a active) Equal(other active) bool {
	return a.obs.Equal(other.obs) && a.key.Equal(other.key)
}

// RecordMethods is initialized by the builtin package
var RecordMethods Methods

var gnRecords = Global.Num("Records")

func (SuRecord) Lookup(method string) Value {
	if m := Lookup(RecordMethods, gnObjects, method); m != nil {
		return m
	}
	return (*SuObject)(nil).Lookup(method)
}

// Packable ---------------------------------------------------------

func (r *SuRecord) Pack(buf *pack.Encoder) {
	r.SuObject.pack(buf, PackRecord)
}

func UnpackRecord(s string) *SuRecord {
	r := NewSuRecord()
	unpackObject(s, &r.SuObject)
	return r
}

// old format

func UnpackRecordOld(s string) *SuRecord {
	r := NewSuRecord()
	unpackObjectOld(s, &r.SuObject)
	return r
}
