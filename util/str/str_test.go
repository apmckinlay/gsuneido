package str

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestCapitalized(t *testing.T) {
	Assert(t).That(Capitalized(""), Equals(false))
	Assert(t).That(Capitalized("a"), Equals(false))
	Assert(t).That(Capitalized("abc"), Equals(false))
	Assert(t).That(Capitalized("?"), Equals(false))
	Assert(t).That(Capitalized("A"), Equals(true))
	Assert(t).That(Capitalized("Abc"), Equals(true))
}

func TestCapitalize(t *testing.T) {
	Assert(t).That(Capitalize(""), Equals(""))
	Assert(t).That(Capitalize("#$%"), Equals("#$%"))
	Assert(t).That(Capitalize("abc"), Equals("Abc"))
	Assert(t).That(Capitalize("a"), Equals("A"))
	Assert(t).That(Capitalize("abC"), Equals("AbC"))
}

func TestUnCapitalize(t *testing.T) {
	Assert(t).That(UnCapitalize(""), Equals(""))
	Assert(t).That(UnCapitalize("#$%"), Equals("#$%"))
	Assert(t).That(UnCapitalize("abc"), Equals("abc"))
	Assert(t).That(UnCapitalize("A"), Equals("a"))
	Assert(t).That(UnCapitalize("AbC"), Equals("abC"))
}
