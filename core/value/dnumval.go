// StrVal is a string Value

package value

import (
	"github.com/apmckinlay/gsuneido/util/dnum"
)

type DnumVal dnum.Dnum

func (dn DnumVal) ToInt() int {
	n, _ := dnum.Dnum(dn).Int32()
	return int(n)
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

func (dn DnumVal) Hash2() uint32 {
	return dn.Hash()
}

func (dn DnumVal) Equals(other interface{}) bool {
	return 0 == dnum.Cmp(dnum.Dnum(dn), other.(dnum.Dnum))
}

var _ Value = DnumVal{} // confirm it implements Value
