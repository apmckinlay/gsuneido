// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typealgebra

import "testing"

// every primitive round-trips: what String writes, ParseName reads back
func TestNameRoundTrip(t *testing.T) {
	for i := range primitiveNames {
		p := Primitive(i)
		got, ok := ParseName(p.String())
		if !ok || got != p {
			t.Errorf("%v printed as %q, read back as (%v, %v)", int(p), p.String(), got, ok)
		}
	}
}

func TestParseNameCaseInsensitive(t *testing.T) {
	for _, name := range []string{"string", "String", "STRING"} {
		if p, ok := ParseName(name); !ok || p != TString {
			t.Errorf("ParseName(%q) = (%v, %v), want TString", name, p, ok)
		}
	}
	if _, ok := ParseName("Foo"); ok {
		t.Errorf("ParseName(%q) matched a builtin", "Foo")
	}
	// an alias parses but is not a spelling of its own
	if p, ok := ParseName("record"); !ok || p != TObject {
		t.Errorf("ParseName(record) = (%v, %v), want TObject", p, ok)
	}
}

func TestCanonicalName(t *testing.T) {
	cases := map[string]string{
		"String": "string",
		"FALSE":  "false",
		"object": "object",
		"Foo":    "Foo",    // a class keeps its capitalization
		"record": "record", // an alias is left as the author wrote it
	}
	for in, want := range cases {
		if got, _ := CanonicalName(in); got != want {
			t.Errorf("CanonicalName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCanonicalAnnotation(t *testing.T) {
	cases := map[string]string{
		"String | False": "string|false",
		"number|Foo":     "number|Foo",
		"string":         "string",
	}
	for in, want := range cases {
		if got := CanonicalAnnotation(in); got != want {
			t.Errorf("CanonicalAnnotation(%q) = %q, want %q", in, got, want)
		}
	}
}

// what NativeAnnotation writes must parse back to the type it was given
func TestNativeAnnotationRoundTrip(t *testing.T) {
	for _, ty := range []DynType{
		TString, TNumber, TFalse, TObject,
		U(TString, TFalse),
		U(Instance{Class: "Foo"}, TNumber),
		Instance{Class: "Point"},
	} {
		s, ok := NativeAnnotation(ty)
		if !ok {
			t.Errorf("%v is not expressible", ty)
			continue
		}
		// Union holds a slice, so compare the rendering, not the value
		back, err := ParseAnnotation(s)
		if err != nil || back.String() != ty.String() {
			t.Errorf("%v wrote %q, read back as %v (err %v)", ty, s, back, err)
		}
	}
}

// a dirty union prints, but cannot be declared
func TestNativeAnnotationRejectsDirty(t *testing.T) {
	dirty := MarkDirty(TString)
	if _, ok := NativeAnnotation(dirty); ok {
		t.Errorf("dirty union reported as expressible")
	}
	if got := dirty.String(); got != "?|string" {
		t.Errorf("dirty union printed as %q", got)
	}
}

// two nominal arms fold in the order they were given, so only the sort in
// joinArms makes the printed union canonical (found by TestPropStringCanonical)
func TestUnionOfTwoInstancesPrintsCanonically(t *testing.T) {
	ca, cb := Instance{Class: "Ca"}, Instance{Class: "Cb"}
	ab, ba := U(U(ca, cb), TUnknown), U(U(cb, ca), TUnknown)
	if ab.String() != ba.String() {
		t.Errorf("non-canonical: %q vs %q", ab, ba)
	}
	if ab.String() != "?|Ca|Cb" {
		t.Errorf("got %q", ab)
	}
}

func TestAnnotationSpan(t *testing.T) {
	cases := []struct {
		src  string
		want string // the span's text, or "" when there is none
	}{
		{"x: number = 0", ": number"},
		{"x :string|false)", " :string|false"},
		{"x : String | False)", " : String | False"},
		{"x = 0", ""},
		{"x:\n\tnumber", ""}, // not on one line: left alone
		{"x:", ""},           // a colon with no name
	}
	for _, c := range cases {
		end, ok := AnnotationSpan(c.src, 1)
		got := ""
		if ok {
			got = c.src[1:end]
		}
		if got != c.want {
			t.Errorf("AnnotationSpan(%q) = %q, want %q", c.src, got, c.want)
		}
	}
}
