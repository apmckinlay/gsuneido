// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"fmt"
	"slices"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
)

const codeFolderLimit = 400

func codeFoldersTool(library, path string) (codeFoldersOutput, error) {
	th := core.NewThread(core.MainThread)
	defer th.Close()
	libs := th.Dbms().Libraries()
	if !slices.Contains(libs, library) {
		return codeFoldersOutput{}, fmt.Errorf("library not found: %s", library)
	}
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	st := core.NewSuTran(tran, false)

	normalized := normalizeFolderPath(path)
	parent, err := resolveFolderParent(th, tran, library, normalized)
	if err != nil {
		return codeFoldersOutput{}, err
	}

	children, err := codeFolderChildren(th, tran, st, library, parent)
	if err != nil {
		return codeFoldersOutput{}, err
	}
	return codeFoldersOutput{
		Library:  library,
		Path:     normalized,
		Children: children,
	}, nil
}

func normalizeFolderPath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "/")
	return path
}

func resolveFolderParent(th *core.Thread, tran core.ITran, library, path string) (int, error) {
	if path == "" {
		return 0, nil
	}
	segments := strings.Split(path, "/")
	parent := 0
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		num, err := folderNum(th, tran, library, parent, segment)
		if err != nil {
			return 0, err
		}
		parent = num
	}
	return parent, nil
}

func folderNum(th *core.Thread, tran core.ITran, library string, parent int, name string) (int, error) {
	folderArgs := core.SuObjectOf(core.SuStr(library))
	folderArgs.Set(core.SuStr("group"), core.IntVal(parent))
	folderArgs.Set(core.SuStr("name"), core.SuStr(name))
	row, hdr, _ := tran.Get(th, folderArgs, core.Only)
	if row == nil {
		leafArgs := core.SuObjectOf(core.SuStr(library))
		leafArgs.Set(core.SuStr("group"), core.IntVal(-1))
		leafArgs.Set(core.SuStr("parent"), core.IntVal(parent))
		leafArgs.Set(core.SuStr("name"), core.SuStr(name))
		leaf, _, _ := tran.Get(th, leafArgs, core.Only)
		if leaf != nil {
			return 0, fmt.Errorf("path segment is not a folder: %s", name)
		}
		return 0, fmt.Errorf("path not found: %s", name)
	}
	num, err := intValue(row.GetVal(hdr, "num", th, nil), "num")
	if err != nil {
		return 0, err
	}
	return num, nil
}

func codeFolderChildren(th *core.Thread, tran core.ITran, st *core.SuTran, library string, parent int) (children []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("code folders query failed: %v", r)
		}
	}()
	query := fmt.Sprintf("%s where parent = %d and group >= -1 sort name", library, parent)
	q := tran.Query(query, nil)
	hdr := q.Header()
	for row, _ := q.Get(th, core.Next); row != nil; row, _ = q.Get(th, core.Next) {
		if len(children) >= codeFolderLimit {
			children = append(children, "...")
			break
		}
		nameVal := row.GetVal(hdr, "name", th, st)
		name, ok := nameVal.ToStr()
		if !ok {
			continue
		}
		group, err := intValue(row.GetVal(hdr, "group", th, st), "group")
		if err != nil {
			return nil, err
		}
		if group > -1 {
			name += "/"
		}
		children = append(children, name)
	}
	return children, nil
}

func intValue(val core.Value, field string) (int, error) {
	if val == nil {
		return 0, fmt.Errorf("%s column not found or null", field)
	}
	if i, ok := val.ToInt(); ok {
		return i, nil
	}
	return 0, fmt.Errorf("%s column is not an integer", field)
}
