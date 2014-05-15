package verify

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestVerify(t *testing.T) {
	That(true) // does nothing
	Assert(t).That(func() { That(false) }, Panics("verify failed"))
}
