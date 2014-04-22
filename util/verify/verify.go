package verify

func Verify(cond bool) {
	if !cond {
		panic("verify failed")
	}
}
