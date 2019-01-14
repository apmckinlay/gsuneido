package runtime

// TODO default, rules, observers, etc.

// SuRecord is an SuObject with observers and rules
type SuRecord struct {
	SuObject
}

func NewSuRecord() *SuRecord {
	return &SuRecord{SuObject{defval: EmptyStr}}
}

func (*SuRecord) TypeName() string {
	return "Record"
}

func (r *SuRecord) String() string {
	s := r.SuObject.String()
	return "[" + s[2:len(s) - 1] + "]"
}

func (r *SuRecord) Show() string {
	s := r.SuObject.Show()
	return "[" + s[2:len(s) - 1] + "]"
}

func (r *SuRecord) Pack(buf []byte) []byte {
	return r.SuObject.pack(buf, packRecord)
}

func UnpackRecord(buf []byte) *SuRecord {
	r := NewSuRecord()
	unpackObject(buf, &r.SuObject)
	return r
}

// TODO override SuObject Lookup
