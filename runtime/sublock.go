package runtime

// SuBlock is an instance of a closure block
type SuBlock struct {
	SuFunc
	locals []Value
	this   Value
}

// Value interface

var _ Value = (*SuBlock)(nil)

func (b *SuBlock) String() string {
	return "/* block */"
}

func (b *SuBlock) Call(t *Thread, as *ArgSpec) Value {
	bf := &b.SuFunc

	// normally done by SuFunc Call
	args := t.Args(&b.ParamSpec, as)

	// copy args
	for i := 0; i < int(b.Nparams); i++ {
		b.locals[int(bf.Offset)+i] = args[i]
	}

	// normally done by Thread Call
	t.frames[t.fp] = Frame{fn: bf, locals: b.locals, this: b.this}
	defer func(fp int) { t.fp = fp }(t.fp)
	t.fp++
	return t.Run()
}

// TypeName returns the Suneido name for the type (Value interface)
func (*SuBlock) TypeName() string {
	return "Block"
}
