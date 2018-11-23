package runtime

// TODO default, rules, observers, etc.

// SuRecord is an SuObject with observers, rules, and default values
type SuRecord struct {
	SuObject
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
