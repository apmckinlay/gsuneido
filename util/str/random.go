package str

import (
	"math/rand"
	"strings"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

const alpha = "abcdefghijklmnopqrstuvwxyz"

func Random(min, max int) string {
	return RandomOf(min, max, alpha)
}

func RandomOf(min, max int, chars string) string {
	return randomOf(min, max, chars, rand.Intn)
}

func randomOf(min, max int, chars string, randIntn func(int) int) string {
	n := min + randIntn(1+max-min)
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < n; i++ {
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
	randIntn := rand.Intn
	if len(seed) > 0 {
		randIntn = rand.New(rand.NewSource(seed[0])).Intn
	}
	return func() string {
		var key string
		for i := 0; i < 10; i++ {
			key = randomOf(min, max, chars, randIntn)
			if _, ok := prev[key]; !ok {
				prev[key] = mark
				return key
			}
		}
		panic("str.UniqueRandomOf too many duplicates")
	}
}
