package str

import (
	"math/rand"
	"strings"
)

const chars = "abcdefghijklmnopqrstuvwxyz"

func Random(min, max int) string {
	n := min + rand.Intn(1+max-min)
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < n; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
