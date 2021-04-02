// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/meta"
)

func DoRequest(db *db19.Database, request string) {
	r := ParseRequest(request)
	r.execute(db)
}

//-------------------------------------------------------------------

type createRequest struct {
	schema Schema
}

func (r *createRequest) String() string {
	return "create " + r.schema.String()
}

func (r *createRequest) execute(db *db19.Database) {
	ts := &meta.Schema{Schema: r.schema}
	ts.Ixspecs()
	ov := make([]*index.Overlay, len(ts.Indexes))
	for i := range ov {
		bt := btree.CreateBtree(db.Store, nil)
		ov[i] = index.OverlayFor(bt)
	}
	ti := &meta.Info{Table: r.schema.Table, Indexes: ov}
	db.LoadedTable(ts, ti)
}

//-------------------------------------------------------------------

type ensureRequest struct {
	schema Schema
}

func (r *ensureRequest) String() string {
	return "ensure " + r.schema.String()
}

func (r *ensureRequest) execute(db *db19.Database) {
	//TODO
}

//-------------------------------------------------------------------

type renameRequest struct {
	from string
	to   string
}

func (r *renameRequest) String() string {
	return "rename " + r.from + " to " + r.to
}

func (r *renameRequest) execute(db *db19.Database) {
	//TODO
}

//-------------------------------------------------------------------

type alterRenameRequest struct {
	table string
	from  []string
	to    []string
}

func (r *alterRenameRequest) String() string {
	s := "alter " + r.table + " rename "
	sep := ""
	for i, from := range r.from {
		s += sep + from + " to " + r.to[i]
		sep = ", "
	}
	return s
}

func (r *alterRenameRequest) execute(db *db19.Database) {
	//TODO
}

//-------------------------------------------------------------------

type alterCreateRequest struct {
	schema Schema
}

func (r *alterCreateRequest) String() string {
	return "alter " + strings.Replace(r.schema.String(), " ", " create ", 1)
}

func (r *alterCreateRequest) execute(db *db19.Database) {
	//TODO
}

//-------------------------------------------------------------------

type alterDropRequest struct {
	schema Schema
}

func (r *alterDropRequest) String() string {
	return "alter " + strings.Replace(r.schema.String(), " ", " drop ", 1)
}

func (r *alterDropRequest) execute(db *db19.Database) {
	//TODO
}

//-------------------------------------------------------------------

type dropRequest struct {
	table string
}

func (r *dropRequest) String() string {
	return "drop " + r.table
}

func (r *dropRequest) execute(db *db19.Database) {
	//TODO
}
