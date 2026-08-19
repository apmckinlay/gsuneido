// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typealgebra

import (
	"fmt"
	"sort"
	"strings"
)

const armSep = "|"
const dirtyArm = "?"

var primitiveNames = [...]string{
	"unknown", "void", "boolean", "false", "true", "number", "string", "date",
	"function", "block", "class", "object", "sequence",
}

var primitiveByName = func() map[string]Primitive {
	m := make(map[string]Primitive, len(primitiveNames))
	for i, n := range primitiveNames {
		m[n] = Primitive(i)
	}
	return m
}()

var typeAliases = map[string]Primitive{
	"record": TObject,
}

func ParseName(name string) (Primitive, bool) {
	lower := strings.ToLower(name)
	if p, ok := primitiveByName[lower]; ok {
		return p, true
	}
	p, ok := typeAliases[lower]
	return p, ok
}

func CanonicalName(name string) (string, bool) {
	if p, ok := primitiveByName[strings.ToLower(name)]; ok {
		return p.String(), true
	}
	return name, false
}

func ParseAnnotation(s string) (DynType, error) {
	if s == "" {
		return TUnknown, nil
	}
	var result DynType
	var errs []string
	for raw := range strings.SplitSeq(s, armSep) {
		name := strings.TrimSpace(raw)
		if name == "" {
			errs = append(errs, "empty type alternative")
			continue
		}
		t, err := parseArm(name)
		if err != nil {
			errs = append(errs, err.Error())
			t = TUnknown
		}
		if result == nil {
			result = t
		} else {
			result = U(result, t)
		}
	}
	if result == nil {
		result = TUnknown
	}
	if len(errs) > 0 {
		return result, fmt.Errorf("annotation %q: %s", s, strings.Join(errs, "; "))
	}
	return result, nil
}

func parseArm(name string) (DynType, error) {
	if t, ok := ParseName(name); ok {
		return t, nil
	}
	if name[0] >= 'A' && name[0] <= 'Z' {
		return Instance{Class: name}, nil
	}
	return TUnknown, fmt.Errorf("unknown type name %q", name)
}

func NativeAnnotation(ty DynType) (string, bool) {
	switch t := ty.(type) {
	case Primitive:
		if t == TUnknown {
			return "", false
		}
		return t.String(), true
	case Instance:
		if t.Class == "" {
			return "", false
		}
		name, _ := CanonicalName(t.Class)
		return name, true
	case Union:
		if t.IsDirty {
			return "", false // the ? alternative is not expressible
		}
		parts := make([]string, 0, len(t.Types))
		for _, m := range t.Types {
			s, ok := NativeAnnotation(m)
			if !ok {
				return "", false
			}
			parts = append(parts, s)
		}
		if len(parts) == 0 {
			return "", false
		}
		return joinArms(parts), true
	default:
		return "", false
	}
}

func CanonicalAnnotation(decl string) string {
	parts := strings.Split(decl, armSep)
	for i, p := range parts {
		name, _ := CanonicalName(strings.TrimSpace(p))
		parts[i] = name
	}
	return strings.Join(parts, armSep)
}

func joinArms(parts []string) string {
	sorted := append([]string(nil), parts...)
	sort.Strings(sorted)
	return strings.Join(sorted, armSep)
}

func AnnotationSpan(src string, from int) (int, bool) {
	i := skipBlank(src, from)
	if i >= len(src) || src[i] != ':' {
		return 0, false
	}
	i++
	for {
		i = skipBlank(src, i)
		j := i
		for j < len(src) && isIdentChar(src[j]) {
			j++
		}
		if j == i {
			return 0, false // a colon with no type name after it
		}
		i = j
		if k := skipBlank(src, i); k < len(src) && src[k] == armSep[0] {
			i = k + 1
			continue
		}
		return i, true
	}
}

func skipBlank(src string, i int) int {
	for i < len(src) && (src[i] == ' ' || src[i] == '\t') {
		i++
	}
	return i
}

func isIdentChar(c byte) bool {
	return c == '_' || c == '?' || c == '!' ||
		('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9')
}
