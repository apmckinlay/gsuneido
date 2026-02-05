// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
)

const (
	isoDate     = "yyyy-MM-dd"
	isoDateTime = "yyyy-MM-dd'T'HH:mm:ss.SSS"
)

func jsonValue(v core.Value, depth int) any {
	if depth > 5 {
		return "<max depth>"
	}
	if v == nil {
		return nil
	}
	switch x := v.(type) {
	case core.SuBool:
		return bool(x)
	case core.SuTimestamp:
		// SuTimestamp has extra precision beyond milliseconds.
		// Represent it as ISO 8601 with a 6-digit fractional second.
		// e.g. 2026-02-03T12:34:56.123007
		s := x.String()
		if len(s) >= 3 {
			return x.Format(isoDateTime) + s[len(s)-3:]
		}
		return x.Format(isoDateTime)
	case core.SuDate:
		if strings.Contains(x.String(), ".") {
			return x.Format(isoDateTime)
		}
		return x.Format(isoDate)
	default:
		if v.Type() == types.Number {
			if n := jsonNumber(v); n != nil {
				return n
			}
		}
		if c, ok := v.ToContainer(); ok {
			if c.NamedSize() == 0 {
				list := make([]any, 0, c.ListSize())
				for i := range c.ListSize() {
					list = append(list, jsonValue(c.ListGet(i), depth+1))
				}
				return list
			}
			m := make(map[string]any, c.NamedSize()+c.ListSize())
			iter := c.Iter2(true, true)
			for k, v := iter(); v != nil; k, v = iter() {
				// TODO: warn if converting keys to strings causes collisions
				// (Suneido keys can be any type)
				m[jsonKey(k, depth+1)] = jsonValue(v, depth+1)
			}
			return m
		}
		if s, ok := v.ToStr(); ok {
			return s
		}
		return v.String()
	}
}

func jsonNumber(v core.Value) any {
	s, ok := v.AsStr()
	if !ok {
		return nil
	}
	if s == "inf" || s == "-inf" {
		return s
	}
	if strings.HasPrefix(s, "-.") {
		s = "-0." + s[2:]
	} else if strings.HasPrefix(s, ".") {
		s = "0" + s
	}
	if !json.Valid([]byte(s)) {
		return nil
	}
	return json.RawMessage([]byte(s))
}

func jsonKey(v core.Value, depth int) string {
	key := jsonValue(v, depth)
	if key == nil {
		return "null"
	}
	if s, ok := key.(string); ok {
		return s
	}
	b, err := json.Marshal(key)
	if err != nil {
		return fmt.Sprint(key)
	}
	return string(b)
}
