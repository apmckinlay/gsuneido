// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typealgebra

type DynType interface {
	String() string
	isDynType()
}

type Primitive int

const (
	TUnknown Primitive = iota
	TVoid              // bottom type, nothing inhabits TVoid

	TBoolean
	TFalse // TFalse <: TBoolean
	TTrue  // TTrue <: TBoolean

	TNumber
	TString
	TDate

	TFunction
	TBlock
	TClass

	TObject
	TSequence // for typechecking purposes we say sequences are subset of object (though suneido runtime semantics differ a bit)

)

type Instance struct {
	Class string
}

func (i Instance) isDynType() {}

func (i Instance) String() string {
	name, _ := CanonicalName(i.Class)
	return name
}

type Union struct {
	Types   []DynType
	IsDirty bool
}

func (u Union) isDynType() {}

func U(a, b DynType) DynType {
	ua, aIsUnion := a.(Union)
	ub, bIsUnion := b.(Union)
	// propagate dirty: TUnknown operands or dirty union operands make the result dirty
	dirty := a == TUnknown || b == TUnknown ||
		(aIsUnion && ua.IsDirty) || (bIsUnion && ub.IsDirty)
	var types []DynType
	switch {
	case aIsUnion && bIsUnion:
		types = append(append([]DynType{}, ua.Types...), ub.Types...)
	case aIsUnion:
		types = append(append([]DynType{}, ua.Types...), b)
	case bIsUnion:
		types = append(append([]DynType{}, ub.Types...), a)
	default:
		types = []DynType{a, b}
	}
	return Union{Types: types, IsDirty: dirty}.Fold()
}

func MarkDirty(t DynType) DynType { return U(t, TUnknown) }

// Equal compares DynTypes structurally. Union holds a slice, so it is not
// comparable - `==` on two Union values is a runtime panic. Never use `==` on
// a DynType, use this.
func Equal(a, b DynType) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	ua, aIsU := a.(Union)
	ub, bIsU := b.(Union)
	if aIsU != bIsU {
		return false
	}
	if !aIsU {
		return a == b // neither is a Union, so both are comparable
	}
	if len(ua.Types) != len(ub.Types) || ua.IsDirty != ub.IsDirty {
		return false
	}
	for _, ta := range ua.Types {
		found := false
		for _, tb := range ub.Types {
			if Equal(ta, tb) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (p Primitive) isDynType() {}

func (p Primitive) String() string {
	return primitiveNames[p]
}

func (u Union) Fold() DynType {
	// carry the dirty flag: either pre-set or triggered by a TUnknown in types
	dirty := u.IsDirty
	seen := map[Primitive]bool{}
	var out []DynType

	for _, t := range u.Types {
		if t == TUnknown {
			dirty = true // mark dirty but keep collecting concrete types
			continue
		}
		p, isPrim := t.(Primitive)
		if isPrim {
			seen[p] = true
		} else if !containsSameInstance(out, t) {
			out = append(out, t)
		}
	}

	// Collapse boolean primitives
	var singleBool Primitive
	boolCount := 0
	for _, p := range []Primitive{TTrue, TFalse, TBoolean} {
		if seen[p] {
			singleBool = p
			boolCount++
			delete(seen, p)
		}
	}
	switch {
	case boolCount == 1:
		out = append(out, singleBool)
	case boolCount > 1:
		out = append(out, TBoolean)
	}

	for p := Primitive(0); int(p) < len(primitiveNames); p++ {
		if seen[p] {
			out = append(out, p)
		}
	}

	switch len(out) {
	case 0:
		return TUnknown // no concrete types at all
	case 1:
		if dirty {
			return Union{Types: out, IsDirty: true}
		}
		return out[0]
	default:
		return Union{Types: out, IsDirty: dirty}
	}
}

func containsSameInstance(out []DynType, t DynType) bool {
	inst, ok := t.(Instance)
	if !ok {
		return false
	}
	for _, o := range out {
		if oi, ok := o.(Instance); ok && oi == inst {
			return true
		}
	}
	return false
}

func (u Union) Contains(want DynType) bool {
	for _, t := range u.Types {
		if Equal(t, want) {
			return true
		}
	}
	return false
}

func (u Union) String() string {
	parts := make([]string, 0, len(u.Types))
	for _, t := range u.Types {
		parts = append(parts, t.String())
	}
	if u.IsDirty {
		parts = append(parts, dirtyArm) // sorts last: ? is not an arm name
	}
	return joinArms(parts)
}
