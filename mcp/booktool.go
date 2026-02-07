// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
)

func bookTool(book, path string) (readBookOutput, error) {
	if path == "/" {
		path = ""
	} else if path != "" && path[0] != '/' {
		path = "/" + path
	}
	if isResPath(path) {
		return readBookOutput{
			Book:     book,
			Path:     path,
			Text:     "",
			Children: []string{},
		}, nil
	}
	th := core.NewThread(core.MainThread)
	defer th.Close()
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	st := core.NewSuTran(tran, false)

	text, err := bookText(th, tran, st, book, path)
	if err != nil {
		return readBookOutput{}, err
	}
	children, err := bookChildren(th, tran, st, book, path)
	if err != nil {
		return readBookOutput{}, err
	}
	return readBookOutput{
		Book:     book,
		Path:     path,
		Text:     text,
		Children: children,
	}, nil
}

func bookText(th *core.Thread, tran core.ITran, st *core.SuTran,
	book, path string) (text string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("book text query failed: %v", r)
		}
	}()
	if path == "" {
		return "", nil
	}
	dir, name := splitPath(path)
	query := fmt.Sprintf("%s where path = %q and name = %q",
		book, dir, name)
	q := tran.Query(query, nil)
	hdr := q.Header()
	row, _ := q.Get(th, core.Next)
	if row == nil {
		return "", nil
	}
	val := row.GetVal(hdr, "text", th, st)
	if val == nil {
		return "", nil
	}
	s, ok := val.ToStr()
	if !ok {
		return "", nil
	}
	return s, nil
}

func bookChildren(th *core.Thread, tran core.ITran, st *core.SuTran,
	book, path string) (children []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("book children query failed: %v", r)
		}
	}()
	children = []string{}
	query := fmt.Sprintf("%s where path = %q sort order, name",
		book, path)
	q := tran.Query(query, nil)
	hdr := q.Header()
	for row, _ := q.Get(th, core.Next); row != nil; row, _ = q.Get(th, core.Next) {
		val := row.GetVal(hdr, "name", th, st)
		if val == nil {
			continue
		}
		if s, ok := val.ToStr(); ok {
			if path == "" && s == "res" {
				continue
			}
			children = append(children, s)
		}
	}
	return children, nil
}

func isResPath(path string) bool {
	return path == "/res" || strings.HasPrefix(path, "/res/")
}

func splitPath(path string) (dir, name string) {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[:i], path[i+1:]
	}
	return "", path
}
