package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// TODO default, rules, observers, etc.
// TODO lazy unpacking of Tuple

// SuRecord is an SuObject with observers and rules
type SuRecord struct {
	SuObject
}

func NewSuRecord() *SuRecord {
	return &SuRecord{SuObject{defval: EmptyStr}}
}

func (r *SuRecord) Copy() *SuRecord {
	return &SuRecord{*r.SuObject.Copy()}
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

// RecordMethods is initialized by the builtin package
var RecordMethods Methods

var gnRecords = Global.Num("Records")

var anSuObject = SuObject{}

func (SuRecord) Lookup(method string) Value {
	if m := Lookup(RecordMethods, gnObjects, method); m != nil {
		return m
	}
	return anSuObject.Lookup(method)
}

// Packable ---------------------------------------------------------

func (r *SuRecord) Pack(buf *pack.Encoder) {
	r.SuObject.pack(buf, packRecord)
}

func UnpackRecord(s string) *SuRecord {
	r := NewSuRecord()
	unpackObject(s, &r.SuObject)
	return r
}
