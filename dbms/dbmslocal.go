// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"

	"slices"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/tools"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/generic/atomics"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

// DbmsLocal implements the Dbms interface using a local database
// i.e. standalone
type DbmsLocal struct {
	db        *db19.Database
	libraries atomics.Value[[]string]
	badlibs   atomic.Bool // limits logging
}

func NewDbmsLocal(db *db19.Database) *DbmsLocal {
	dbms := DbmsLocal{db: db}
	dbms.libraries.Store([]string{"stdlib"})
	return &dbms
}

// Dbms interface

var _ IDbms = (*DbmsLocal)(nil)

func (dbms *DbmsLocal) Admin(admin string, sv *Sviews) {
	trace.Dbms.Println("Admin", admin)
	qry.DoAdmin(dbms.db, admin, sv)
}

func (dbms *DbmsLocal) Auth(th *Thread, s string) bool {
	if DbmsAuth {
		panic("already authorized")
	}
	if !auth(th, s) {
		return false
	}
	DbmsAuth = true
	th.SetDbms(dbms) // not strictly necessary, removes unauth wrap
	return true
}

func auth(th *Thread, s string) bool {
	if AuthUser(th, s, th.Nonce) {
		th.Nonce = ""
		return true
	}
	defer th.Suneido.Store(th.Suneido.Load())
	th.Suneido.Store(nil) // use main Suneido object
	return AuthToken(s)
}

func (dbms *DbmsLocal) Check() string {
	if err := dbms.db.Check(); err != nil {
		return err.Error()
	}
	return ""
}

func (*DbmsLocal) Connections() Value {
	if options.Action == "server" {
		return connections()
	}
	return &SuObject{}
}

func (dbms *DbmsLocal) Cursor(query string, sv *Sviews) ICursor {
	tran := dbms.db.NewReadTran()
	q, fixcost, varcost := buildQuery(query, tran, sv, qry.CursorMode)
	trace.Query.Println("cursor", fixcost+varcost, "-", query)
	return cursorLocal{queryLocal{
		Query: q, cost: fixcost + varcost, mode: qry.CursorMode}}
}

func buildQuery(query string, tran qry.QueryTran, sv *Sviews,
	mode qry.Mode) (qry.Query, int, int) {
	q := qry.ParseQuery(query, tran, sv)
	q, fixcost, varcost := qry.Setup(q, mode, tran)
	qry.Warnings(query, q)
	return q, fixcost, varcost
}

func (*DbmsLocal) Cursors() int {
	return 0
}

func (dbms *DbmsLocal) Corrupted() bool {
	return dbms.db.IsCorrupted()
}

func (dbms *DbmsLocal) DisableTrigger(table string) {
	dbms.db.DisableTrigger(table)
}
func (dbms *DbmsLocal) EnableTrigger(table string) {
	dbms.db.EnableTrigger(table)
}

func (dbms *DbmsLocal) Dump(table, to, publicKey string) string {
	var err error
	if table == "" {
		if to == "" {
			to = "database.su"
		}
		_, _, err = tools.Dump(dbms.db, to, publicKey)
	} else {
		if to == "" {
			to = table + ".su"
		}
		_, err = tools.DumpDbTable(dbms.db, table, to, publicKey)
	}
	if err != nil {
		return err.Error()
	}
	return ""
}

func (*DbmsLocal) Exec(th *Thread, v Value) Value {
	defer th.Suneido.Store(th.Suneido.Load())
	th.Suneido.Store(nil) // use main Suneido object
	trace.Dbms.Println("Exec", v)
	fname := ToStr(ToContainer(v).ListGet(0))
	if i := strings.IndexByte(fname, '.'); i != -1 {
		ob := Global.GetName(th, fname[:i])
		m := fname[i+1:]
		return th.CallLookupEach1(ob, m, v)
	}
	fn := Global.GetName(th, fname)
	return th.CallEach1(fn, v)
}

func (dbms *DbmsLocal) Final() int {
	return dbms.db.Final()
}

// Get handles QueryFirst, QueryLast, Query1, QueryExists?
func (dbms *DbmsLocal) Get(
	th *Thread, query Value, dir Dir) (Row, *Header, string) {
	tran := dbms.db.NewReadTran()
	defer tran.Complete()
	return get(th, tran, query, dir)
}

func (dbms *DbmsLocal) Info() Value {
	ob := &SuObject{}
	ob.Set(SuStr("currentSize"), Int64Val(int64(dbms.db.Size())))
	ob.Set(SuStr("timeoutMin"), IntVal(int(options.TimeoutMinutes)))
	return ob
}

func (*DbmsLocal) Kill(addr string) int {
	if options.Action == "server" {
		return kill(addr)
	}
	return 0
}

func (dbms *DbmsLocal) Load(table, from, privateKey, passphrase string) int {
	if from == "" {
		from = table + ".su"
	}
	n, err := tools.LoadDbTable(table, from, privateKey, passphrase, dbms.db)
	if err != nil {
		panic(err.Error())
	}
	return n
}

// LibGet returns a list of strings.
// The strings are in pairs - "<lib>:<name>[__tag]", "definition".
// The order is significant - first by Libraries() and then by LibraryTags.
// Later definitions can override or inherit from earlier ones.
// NOTE: libload depends on every library being in the results.
func (dbms *DbmsLocal) LibGet(name string) []string {
	defer func() {
		if e := recover(); e != nil {
			// dbg.PrintStack()
			panic("error loading " + name + " " + fmt.Sprint(e))
		}
	}()

	defs := make([]string, 0, 4)
	rt := dbms.db.NewReadTran()
	libs := dbms.libraries.Load()
	for _, lib := range libs {
		defs = dbms.LibGet1(rt, lib, name, defs)
	}
	return defs
}

var libKey = []string{"name", "group"} // const

// LibGet1 returns the definition(s) for the given library and name.
// There may be multiple due to tags.
func (dbms *DbmsLocal) LibGet1(rt *db19.ReadTran, lib, name string, defs []string) []string {
	var tag string
	defer func() {
		if e := recover(); e != nil {
			log.Println("libGet", lib, name, tag, e)
		}
	}()
	ix := rt.GetIndex(lib, libKey)
	if ix == nil {
		dbms.liblog(lib)
		return append(defs, lib, "")
	}
	fld := rt.ColToFld(lib, "text")
	if fld == -1 {
		dbms.liblog(lib)
		return append(defs, lib, "")
	}
	found := false
	var rb ixkey.Encoder
	for _, tag := range options.LibraryTags {
		nametag := name + tag
		rb.Add(Pack(SuStr(nametag)))
		rb.Add(Pack(SuInt(-1)))
		key := rb.String()
		off := ix.Lookup(key)
		if off != 0 {
			defs = append(defs, lib+tag, rt.GetRecord(off).GetStr(fld))
			found = true
		}
	}
	if !found {
		defs = append(defs, lib, "")
	}
	return defs
}

func (dbms *DbmsLocal) liblog(lib string) {
	if !dbms.badlibs.Swap(true) {
		log.Println("ERROR: invalid library: " + lib)
	}
}

func (dbms *DbmsLocal) Libraries() []string {
	// library list is not mutated so it's thread safe to return
	return dbms.libraries.Load()
}

func (*DbmsLocal) Log(s string) {
	log.Println(s)
}

func (*DbmsLocal) Nonce(th *Thread) string {
	th.Nonce = Nonce()
	return th.Nonce
}

func (*DbmsLocal) Run(th *Thread, s string) Value {
	defer th.Suneido.Store(th.Suneido.Load())
	th.Suneido.Store(nil) // use main Suneido object
	trace.Dbms.Println("Run", s)
	return compile.EvalString(th, s)
}

func (dbms *DbmsLocal) Schema(table string) string {
	return dbms.db.Schema(table)
}

func (*DbmsLocal) SessionId(th *Thread, id string) string {
	if id != "" {
		th.SetSession(id)
	}
	return th.Session()
}

func (dbms *DbmsLocal) Size() uint64 {
	return dbms.db.Size()
}

func (*DbmsLocal) Timestamp() SuDate {
	return db19.Timestamp()
}

func (*DbmsLocal) Token() string {
	return Token()
}

func (dbms *DbmsLocal) Transaction(update bool) ITran {
	if update {
		if t := dbms.db.NewUpdateTran(); t != nil {
			return &UpdateTranLocal{UpdateTran: t}
		}
		panic(fmt.Sprintf("too many overlapping update transactions (%d)",
			db19.MaxTrans))
	}
	return &ReadTranLocal{ReadTran: dbms.db.NewReadTran()}
}

// Transactions only returns the update transactions
func (dbms *DbmsLocal) Transactions() *SuObject {
	trans := dbms.db.Transactions()
	if trans == nil {
		return SuObjectOf(Zero) // corrupt
	}
	slices.Sort(trans)
	var ob SuObject
	for _, t := range trans {
		ob.Add(IntVal(t))
	}
	return &ob
}

func (dbms *DbmsLocal) Unuse(lib string) bool {
	return dbms.updateLibraries(func(libs []string) []string {
		if lib == "stdlib" || !slices.Contains(libs, lib) {
			return nil
		}
		return slc.Without(libs, lib) // copy on write
	})
}

func (dbms *DbmsLocal) Use(lib string) bool {
	return dbms.updateLibraries(func(libs []string) []string {
		if slices.Contains(libs, lib) {
			return nil
		}
		dbms.checkLibrary(lib)
		return append(libs, lib)
	})
}

func (dbms *DbmsLocal) checkLibrary(lib string) {
	rt := dbms.db.NewReadTran()
	if rt.GetIndex(lib, libKey) == nil || rt.ColToFld(lib, "text") == -1 {
		panic("Use: invalid library: " + lib)
	}
}

func (dbms *DbmsLocal) updateLibraries(fn func(libs []string) []string) bool {
	oldlibs := dbms.libraries.Load()
	newlibs := fn(oldlibs)
	if newlibs == nil {
		return false
	}
	dbms.badlibs.Store(false) // reset logging
	return slices.Equal(oldlibs, dbms.libraries.Swap(newlibs))
}

func (dbms *DbmsLocal) Unwrap() IDbms {
	return dbms
}

func (dbms *DbmsLocal) FormatQuery(query string) string {
	t := dbms.db.NewReadTran()
	defer t.Complete()
	return qry.Format(t, query)
}

func (dbms *DbmsLocal) Close() {
	dbms.db.Close()
}

// ReadTranLocal --------------------------------------------------------

func init() {
	qry.MakeSuTran = func(qt qry.QueryTran) *SuTran {
		if qt == nil {
			return nil
		}
		switch t := qt.(type) {
		case *ReadTranLocal:
			return NewSuTran(t, false)
		case *UpdateTranLocal:
			return NewSuTran(t, true)
		case *db19.ReadTran:
			return NewSuTran(&ReadTranLocal{ReadTran: t}, false)
		case *db19.UpdateTran:
			return NewSuTran(&UpdateTranLocal{UpdateTran: t}, true)
		}
		panic(fmt.Sprintf("NewSuTran unhandled type %#v", qt))
	}
	db19.MakeSuTran = func(ut *db19.UpdateTran) *SuTran {
		return NewSuTran(&UpdateTranLocal{UpdateTran: ut}, true)
	}
}

type ReadTranLocal struct {
	*db19.ReadTran
}

func (t ReadTranLocal) Get(th *Thread, query Value, dir Dir) (Row, *Header, string) {
	return get(th, t.ReadTran, query, dir)
}

func (t ReadTranLocal) Query(query string, sv *Sviews) IQuery {
	q, fixcost, varcost := buildQuery(query, t.ReadTran, sv, qry.ReadMode)
	trace.Query.Println(fixcost+varcost, "-", query)
	return queryLocal{Query: q, cost: fixcost + varcost, mode: qry.ReadMode}
}

func (t ReadTranLocal) Action(*Thread, string) int {
	panic("cannot do action in read-only transaction")
}

// UpdateTranLocal --------------------------------------------------------

type UpdateTranLocal struct {
	*db19.UpdateTran
}

func (t UpdateTranLocal) Get(th *Thread, query Value, dir Dir) (Row, *Header, string) {
	return get(th, t.UpdateTran, query, dir)
}

func (t UpdateTranLocal) Query(query string, sv *Sviews) IQuery {
	q, fixcost, varcost := buildQuery(query, t.UpdateTran, sv, qry.UpdateMode)
	trace.Query.Println("update", fixcost+varcost, "-", query)
	return queryLocal{Query: q, cost: fixcost + varcost, mode: qry.UpdateMode}
}

func (t UpdateTranLocal) Action(th *Thread, action string) int {
	defer th.Suneido.Store(th.Suneido.Load())
	th.Suneido.Store(nil) // use main Suneido object
	trace.Dbms.Println("Action", action)
	return qry.DoAction(th, t.UpdateTran, action)
}

func (t UpdateTranLocal) Update(th *Thread, table string, oldoff uint64, newrec Record) uint64 {
	defer th.Suneido.Store(th.Suneido.Load())
	th.Suneido.Store(nil) // use main Suneido object
	trace.Dbms.Println("Update", table)
	return t.UpdateTran.Update(th, table, oldoff, newrec)
}

// queryLocal

type queryLocal struct {
	// Query is embedded so most methods are "inherited" directly
	qry.Query
	keys []string // cache
	cost qry.Cost
	mode qry.Mode
}

func (q queryLocal) Keys() []string {
	if q.keys == nil {
		keys := q.Query.Keys()
		list := make([]string, len(keys))
		for i, k := range keys {
			list[i] = str.Join(",", k)
		}
		q.keys = list
	}
	return q.keys
}

func (q queryLocal) Strategy(formatted bool) string {
	var strategy string
	if formatted {
		strategy = qry.Strategy(q.Query) + "\n"
	} else {
		strategy = qry.String(q.Query) + " "
	}
	n, _ := q.Nrows()
	return fmt.Sprint(strategy,
		"[nrecs~ ", n, " cost~ ", q.cost, " ", q.mode, "]")
}

func (q queryLocal) Order() []string {
	return q.Query.Order()
}

func (q queryLocal) Get(th *Thread, dir Dir) (Row, string) {
	defer th.Suneido.Store(th.Suneido.Load())
	th.Suneido.Store(nil) // use main Suneido object
	row := q.Query.Get(th, dir)
	if row == nil {
		// this is required for SuQuery to stick at eof unidirectionally
		q.Query.Rewind()
	}
	return row, q.Query.Updateable()
}

func (q queryLocal) Tree() Value {
	qry.CalcSelf(q.Query)
	return qry.NewSuQueryNode(q.Query)
}

func (q queryLocal) Close() {
}

// cursorLocal

type cursorLocal struct {
	queryLocal
}

func (q cursorLocal) Get(th *Thread, t ITran, dir Dir) (Row, string) {
	q.Query.SetTran(t.(qry.QueryTran))
	return q.queryLocal.Get(th, dir)
}
