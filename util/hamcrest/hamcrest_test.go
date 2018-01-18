package hamcrest

import gotesting "testing"

func TestHamcrest(t *gotesting.T) {
	assert := Assert(t)
	assert.That(func() { panic("error") }, Panics("error"))
}
