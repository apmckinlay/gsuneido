// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestPackToStr(t *testing.T) {
	assert := assert.T(t)
	assert.This(packToStr("")).Is("''")
	assert.This(packToStr("\xff")).Is("max")
	assert.This(packToStr(string([]byte{PackString}))).Is("PackString")
	assert.This(packToStr(string([]byte{PackMinus}))).Is("PackMinus")
	assert.This(packToStr(string([]byte{PackDate}))).Is("PackDate")
	assert.This(packToStr(string([]byte{PackDate + 1}))).Is("PackDate+1")
	assert.This(packToStr(Pack(SuInt(123)))).Is("123")
	assert.This(packToStr(Pack(SuStr("abc")))).Is("'abc'")
}

func TestPointRange(t *testing.T) {
	pr := pointRange{Org: Pack(SuInt(5))}
	assert.T(t).That(pr.isPoint())
	assert.T(t).That(!pr.isRange())
	assert.This(pr.String()).Is("5")

	pr = pointRange{Org: Pack(SuInt(2)), End: Pack(SuInt(7))}
	assert.T(t).That(!pr.isPoint())
	assert.T(t).That(pr.isRange())
	assert.This(pr.String()).Is("2..7")
	assert.T(t).That(!pr.conflict())

	pr = pointRange{Org: "z", End: "a"}
	assert.T(t).That(pr.conflict())
	assert.This(pr.String()).Is("<empty>")
}

func TestPointRangeIntersect(t *testing.T) {
	point := pointRange{Org: Pack(SuInt(5))}
	assert.This(point.intersect(Pack(SuInt(1)), Pack(SuInt(9)))).Is(point)
	assert.T(t).That(point.intersect(Pack(SuInt(6)), Pack(SuInt(9))).conflict())

	rng := pointRange{Org: Pack(SuInt(2)), End: Pack(SuInt(8))}
	assert.This(rng.intersect(Pack(SuInt(4)), Pack(SuInt(9)))).Is(pointRange{
		Org: Pack(SuInt(4)), End: Pack(SuInt(8)),
	})
	assert.This(rng.intersect(Pack(SuInt(1)), Pack(SuInt(6)))).Is(pointRange{
		Org: Pack(SuInt(2)), End: Pack(SuInt(6)),
	})
	assert.T(t).That(rng.intersect(Pack(SuInt(8)), Pack(SuInt(9))).conflict())
}

func TestIdxSelString(t *testing.T) {
	is := idxSel{
		index:     []string{"a", "b", "c", "d", "e"},
		indexFrac: .25,
		dataFrac:  .33,
		prefixLen: 1,
		prefixRanges: []pointRange{
			{Org: Pack(SuInt(1))},
			{Org: Pack(SuInt(2)), End: Pack(SuInt(4))}},
		skipStart:   1,
		skipLen:     1,
		skipRange:   pointRange{Org: Pack(SuInt(3)), End: Pack(SuInt(6))},
		moreFilters: []string{"d", "e"},
	}
	assert.T(t).This(is.String()).
		Is("(a,b,c,d,e) a: <1 | 2..4> +b: <3..6> d,e = .25 .33")

	is = idxSel{
		index:       []string{"a", "b", "c"},
		indexFrac:   1,
		dataFrac:    .5,
		moreFilters: []string{"c"},
	}
	assert.T(t).This(is.String()).Is("(a,b,c) c = 1 .5")

	is = idxSel{
		index:     []string{"a", "b", "c"},
		indexFrac: .33,
		dataFrac:  .22,
		prefixLen: 2,
		prefixRanges: []pointRange{
			{Org: Pack(SuInt(7))},
		},
	}
	assert.T(t).This(is.String()).Is("(a,b,c) a,b: <7> = .33 .22")

	is = idxSel{
		index:     []string{"a", "b", "c"},
		indexFrac: .2,
		dataFrac:  .2,
		prefixLen: 1,
		prefixRanges: []pointRange{
			{Org: Pack(SuInt(1)), End: ixkey.Max},
		},
		skipStart: 2,
		skipLen:   1,
		skipRange: pointRange{Org: Pack(SuInt(5))},
	}
	assert.T(t).This(is.String()).Is("(a,b,c) a: <1..max> +c: <5> = .2")

	// encoded
	is = idxSel{
		index:     []string{"a", "b"},
		encoded:   true,
		indexFrac: .1,
		dataFrac:  .2,
		prefixLen: 1,
		prefixRanges: []pointRange{{
			Org: ixkey.CompKey(Pack(SuInt(1)), Pack(SuStr("x"))),
			End: ixkey.CompKey(Pack(SuInt(2)), Pack(SuStr("z")))}},
	}
	assert.T(t).This(is.String()).Is("(a,b) a: <1,'x'..2,'z'> = .1 .2")
}

func TestFracStr(t *testing.T) {
	assert.T(t).This(fracStr(.33333)).Is(".33")
	assert.T(t).This(fracStr(.00123)).Is(".0012")
	assert.T(t).This(fracStr(10)).Is("10")
}
