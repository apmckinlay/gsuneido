package runtime

// SuSequence wraps an iterator and instantiates it lazily
type SuSequence struct {
	iter Iter
	// ob is nil until the sequence is instantiated
	ob *SuObject
	CantConvert
}

func NewSuSequence(it Iter) *SuSequence {
	return &SuSequence{iter: it}
}

func (seq *SuSequence) instantiate() {
	if seq.ob != nil {
		return // already instantiated
	}
	if seq.iter.Infinite() {
		panic("can't instantiate infinite sequence")
	}
	seq.ob = &SuObject{}
	for x := seq.iter.Next(); x != nil; x = seq.iter.Next() {
		seq.ob.Add(x)
	}
}

func (seq *SuSequence) Iter() Iter {
	return seq.iter.Dup()
}

// Value interface --------------------------------------------------

var _ Value = (*SuSequence)(nil)

func (seq *SuSequence) String() string {
	if seq.iter.Infinite() {
		return "infiniteSequence"
	}
	seq.instantiate()
	return seq.ob.String()
}

func (seq *SuSequence) ToObject() (*SuObject, bool) {
	seq.instantiate()
	return seq.ob, true
}

func (seq *SuSequence) Get(t *Thread, key Value) Value {
	seq.instantiate()
	return seq.ob.Get(t, key)
}

func (seq *SuSequence) Put(key Value, val Value) {
	seq.instantiate()
	seq.ob.Put(key, val)
}

func (seq *SuSequence) RangeTo(i int, j int) Value {
	seq.instantiate()
	return seq.ob.RangeTo(i, j)
}

func (seq *SuSequence) RangeLen(i int, n int) Value {
	seq.instantiate()
	return seq.ob.RangeLen(i, n)
}

func (seq *SuSequence) Equal(other interface{}) bool {
	seq.instantiate()
	return seq.ob.Equal(other)
}

func (seq *SuSequence) Hash() uint32 {
	seq.instantiate()
	return seq.ob.Hash()
}

func (seq *SuSequence) Hash2() uint32 {
	seq.instantiate()
	return seq.ob.Hash2()
}

func (*SuSequence) TypeName() string {
	return "Object"
}

func (*SuSequence) Order() Ord {
	return ordObject
}

func (seq *SuSequence) Compare(other Value) int {
	seq.instantiate()
	return seq.ob.Compare(other)
}

func (*SuSequence) Call(*Thread, *ArgSpec) Value {
	panic("can't call Object")
}

// SequenceMethods is initialized by the builtin package
var SequenceMethods Methods

func (seq *SuSequence) Lookup(method string) Value {
	if meth := SequenceMethods[method]; meth != nil {
		return meth
	}
	return seq.Lookup(method)
}
