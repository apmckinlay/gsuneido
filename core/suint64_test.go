// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"math"
	rand "math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/pack"
)

func TestPackInt(t *testing.T) {
	//   9223372036854775807
	x := int64(9200000000000000000)
	enc := pack.NewEncoder(100)
	packInt(x, enc)
	p := enc.String()
	y := unpackInt64(p)
	if x != y {
		fmt.Println(x)
		if -9999_9999_9999_9999 < x && x < 9999_9999_9999_9999 {
			fmt.Printf("old %x\n", Pack(SuDnum{Dnum: dnum.FromInt(int64(x))}))
		}
		fmt.Printf("new %x\n", p)
		fmt.Println(y)
		t.FailNow()
	}
}

func TestUnpackNum(t *testing.T) {
	test := func(s string, isint bool) {
		for _, sign := range []string{"", "-"} {
			d := SuDnum{Dnum: dnum.FromStr(sign + s)}
			p := Pack(d)
			v := UnpackNumber(p)
			_, ok := v.(SuDnum)
			assert.Msg(sign + s).This(!ok).Is(isint)
		}
	}
	test("0", true)
	test("1", true)
	test("1.5", false)
	test(".123e2", false)
	test(".123e3", true)
	test(".123e4", true)
	test("1e18", true)
	test("1e19", false)
}

func FuzzIsInt(f *testing.F) {
	f.Add(int8(0), 0)
	f.Add(int8(1), 1)
	f.Add(int8(1), 12)
	f.Add(int8(2), 12)
	f.Fuzz(func(t *testing.T, x int8, coef int) {
		sign := int8(-1)
		if x&1 == 1 {
			sign = -1
		}
		exp := x >> 1
		d := SuDnum{Dnum: dnum.New(sign, uint64(coef), int(exp))}
		// fmt.Println(d)
		s := Pack(d)
		v := UnpackNumber(s)
		assert.This(v).Is(v)
	})
} // go test -fuzz=FuzzIsInt -run=FuzzIsInt ./core

func FuzzPackUnpackInt(f *testing.F) {
	f.Add(0)
	f.Add(1)
	f.Add(-1)
	f.Add(12)
	f.Add(-12)
	f.Add(123)
	f.Add(-123)
	f.Add(1234)
	f.Add(-1234)
	f.Add(math.MaxInt16)
	f.Add(math.MinInt16)
	f.Add(math.MaxInt32)
	f.Add(math.MinInt32)
	f.Add(math.MaxInt64)
	f.Add(math.MinInt64)
	for n := 10; n < 10000000000; n *= 10 {
		f.Add(n)
		f.Add(-n)
		f.Add(n * 12)
		f.Add(n * 123)
	}
	f.Fuzz(func(t *testing.T, i int) {
		x := int64(i)
		packSize := packSizeInt(x)
		enc := pack.NewEncoder(32)
		packInt(x, enc)
		packed := enc.String()
		y := unpackInt64(packed)
		if y != x || packSize != len(packed) || packed[len(packed)-1] == 0 {
			fmt.Println(x)
			fmt.Println("packSize", packSize)
			if -9999_9999_9999_9999 < x && x < 9999_9999_9999_9999 {
				fmt.Printf("old %x\n", Pack(SuDnum{Dnum: dnum.FromInt(int64(x))}))
			}
			fmt.Printf("new %x\n", packed)
			fmt.Println("unpackInt", y)
			t.FailNow()
		}
	})
} // go test -fuzz=FuzzPackUnpackInt -run=FuzzPackUnpackInt ./core

func unpackInt64(s string) int64 {
	if len(s) <= 1 {
		return 0
	}
	sign := int8(+1)
	xor := byte(0)
	if s[0] == PackMinus {
		sign = -1
		xor = 0xff
	}
	exp := int8(s[1] ^ 0x80 ^ xor)
	return int64(ToInt(unpackInt(s, sign, exp, xor)))
}

func BenchmarkPackInt(b *testing.B) {
	for b.Loop() {
		enc := pack.NewEncoder(30)
		n := rand.Int64()
		packInt(n, enc)
		unpackInt64(enc.String())
	}
}
