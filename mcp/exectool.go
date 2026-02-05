// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/str"
)

func exectool(code string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("execute error: %v", r)
		}
	}()

	th := core.NewThread(core.MainThread)
	defer th.Close()

	code = strings.TrimSpace(code)
	src := "function () {\n" + code + "\n}"
	v, warnings := compile.Checked(th, src)
	fn := v.(*core.SuFunc)

	res := th.Call(fn)
	results := []string{}
	if res != nil {
		results = append(results, res.String())
	} else if len(th.ReturnMulti) > 0 {
		for i := len(th.ReturnMulti) - 1; i >= 0; i-- {
			results = append(results, th.ReturnMulti[i].String())
		}
	}
	result = fmt.Sprintf("{\n"+
		"code: %q\n"+
		"warnings: %s\n"+
		"results: %s\n"+
		"}",
		code, quotedCommaList(warnings), str.Join("[, ]", results))
	return
}

func quotedCommaList(ss []string) string {
	if len(ss) == 0 {
		return "[]"
	}
	qs := make([]string, len(ss))
	for i, s := range ss {
		qs[i] = strconv.Quote(s)
	}
	return "[" + strings.Join(qs, ", ") + "]"
}
