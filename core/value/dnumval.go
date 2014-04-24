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

func (dn DnumVal) Get(key Value) Value {
	panic("number does not support get")
}

func (dn DnumVal) Put(key Value, val Value) {
	panic("number does not support put")
}

var _ Value = DnumVal{} // confirm it implements Value
