package str

import (
	"math/rand/v2"
	"strings"
)

// WARNING: not thread safe

var rnd = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

const alpha = "abcdefghijklmnopqrstuvwxyz"

func Random(min, max int) string {
	return RandomOf(min, max, alpha)
}

func RandomOf(min, max int, chars string) string {
	return randomOf(min, max, chars, rand.IntN)
}

func randomOf(min, max int, chars string, randIntn func(int) int) string {
	n := min + randIntn(1+max-min)
	var b strings.Builder
	b.Grow(n)
	for range n {
		b.WriteByte(chars[randIntn(len(chars))])
	}
	return b.String()
}

func UniqueRandom(min, max int, seed ...int64) func() string {
	return UniqueRandomOf(min, max, alpha, seed...)
}

func UniqueRandomOf(min, max int, chars string, seed ...int64) func() string {
	type set struct{}
	var mark set
	prev := map[string]set{}
	randIntN := rand.IntN
	if len(seed) > 0 {
		randIntN = rand.New(rand.NewPCG(uint64(seed[0]), uint64(seed[0]))).IntN
	}
	return func() string {
		var key string
		for range 10 {
			key = randomOf(min, max, chars, randIntN)
			if _, ok := prev[key]; !ok {
				prev[key] = mark
				return key
			}
		}
		panic("str.UniqueRandomOf too many duplicates")
	}
}
