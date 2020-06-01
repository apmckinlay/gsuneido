package str

import (
	"math/rand"
	"strings"
)

func Random(min, max int) string {
	return RandomOf(min, max, "abcdefghijklmnopqrstuvwxyz")
}

func RandomOf(min, max int, chars string) string {
	n := min + rand.Intn(1+max-min)
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < n; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
