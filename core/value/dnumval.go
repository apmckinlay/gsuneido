// StrVal is a string Value

package value

import (
	"github.com/apmckinlay/gsuneido/util/dnum"
)

type DnumVal dnum.Dnum

func (dn DnumVal) ToInt() int32 {
	n, _ := dnum.Dnum(dn).Int32()
	return n
}

func (dn DnumVal) ToDnum() dnum.Dnum {
	return dnum.Dnum(dn)
}

func (dn DnumVal) ToStr() string {
	return dnum.Dnum(dn).String()
}

func (dn DnumVal) String() string {
	return dn.ToStr()
}

func (dn DnumVal) Get(key Value) Value {
	panic("number does not support get")
}

func (dn DnumVal) Put(key Value, val Value) {
	panic("number does not support put")
}

func (dn DnumVal) Hash() uint32 {
	return dnum.Dnum(dn).Hash()
}

func (dn DnumVal) hash2() uint32 {
	return dn.Hash()
}

func (dn DnumVal) Equals(other interface{}) bool {
	if d2, ok := other.(DnumVal); ok {
		return 0 == dnum.Cmp(dnum.Dnum(dn), dnum.Dnum(d2))
	}
	return false
}

var _ Value = DnumVal{} // confirm it implements Value
