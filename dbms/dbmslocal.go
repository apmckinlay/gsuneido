// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"fmt"
	"log"
	"strings"
	"sync"
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
	return &cursorLocal{qcLocal{
		q: q, cost: fixcost + varcost, mode: qry.CursorMode}}
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

// Get implements QueryFirst, QueryLast, Query1
func (dbms *DbmsLocal) Get(
	th *Thread, query string, dir Dir) (Row, *Header, string) {
	tran := dbms.db.NewReadTran()
	defer tran.Complete()
	return get(th, tran, query, dir)
}

func get(th *Thread, tran qry.QueryTran, query string, dir Dir) (Row, *Header, string) {
	defer th.Suneido.Store(th.Suneido.Load())
	th.Suneido.Store(nil) // use main Suneido object
	q := qry.ParseQuery(query, tran, th.Sviews())
	if dir != Only &&
		!strings.Contains(query, "CHECKQUERY SUPPRESS: SORT REQUIRED") {
		if _, ok := q.(*qry.Sort); !ok {
			panic("query first/last require sort")
		}
	}
	q, fixcost, varcost := qry.Setup1(q, qry.ReadMode, tran)
	qry.Warnings(query, q)
	if trace.Query.On() {
		d := map[Dir]string{Only: "one", Next: "first", Prev: "last"}[dir]
		trace.Query.Println(d, fixcost+varcost, "-", query)
	}
	only := false
	if dir == Only {
		only = true
		dir = Next
	}
	row := q.Get(th, dir)
	if row == nil {
		return nil, nil, ""
	}
	if only && !single(q) && q.Get(th, dir) != nil {
		panic("Query1 not unique: " + query)
	}
	return row, q.Header(), q.Updateable()
}

func single(q qry.Query) bool {
	keys := q.Keys()
	return len(keys) == 1 && len(keys[0]) == 0
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

func (dbms *DbmsLocal) LibGet(name string) []string {
	defer func() {
		if e := recover(); e != nil {
			// dbg.PrintStack()
			panic("error loading " + name + " " + fmt.Sprint(e))
		}
	}()

	results := make([]string, 0, 2)
	rt := dbms.db.NewReadTran()
	libs := dbms.libraries.Load()
	for _, lib := range libs {
		s := dbms.LibGet1(rt, lib, name)
		if s != "" {
			results = append(results, lib, string(s))
		}
	}
	return results
}

var libKey = []string{"name", "group"} // const

func (dbms *DbmsLocal) LibGet1(rt *db19.ReadTran, lib, name string) string {
	defer func() {
		if e := recover(); e != nil {
			log.Println("libGet", lib, name, e)
		}
	}()
	ix := rt.GetIndex(lib, libKey)
	if ix == nil {
		dbms.liblog(lib)
		return ""
	}
	fld := rt.ColToFld(lib, "text")
	if fld == -1 {
		dbms.liblog(lib)
		return ""
	}
	var rb ixkey.Encoder
	rb.Add(Pack(SuStr(name)))
	rb.Add(Pack(SuInt(-1)))
	key := rb.String()
	off := ix.Lookup(key)
	if off == 0 {
		return "" // not found
	}
	return rt.GetRecord(off).GetStr(fld)
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

func (t ReadTranLocal) Get(th *Thread, query string, dir Dir) (Row, *Header, string) {
	return get(th, t.ReadTran, query, dir)
}

func (t ReadTranLocal) Query(query string, sv *Sviews) IQuery {
	q, fixcost, varcost := buildQuery(query, t.ReadTran, sv, qry.ReadMode)
	trace.Query.Println(fixcost+varcost, "-", query)
	return &queryLocal{qcLocal: qcLocal{q: q, cost: fixcost + varcost, mode: qry.ReadMode}}
}

func (t ReadTranLocal) Action(*Thread, string) int {
	panic("cannot do action in read-only transaction")
}

// UpdateTranLocal --------------------------------------------------------

type UpdateTranLocal struct {
	*db19.UpdateTran
}

func (t UpdateTranLocal) Get(th *Thread, query string, dir Dir) (Row, *Header, string) {
	return get(th, t.UpdateTran, query, dir)
}

func (t UpdateTranLocal) Query(query string, sv *Sviews) IQuery {
	q, fixcost, varcost := buildQuery(query, t.UpdateTran, sv, qry.UpdateMode)
	trace.Query.Println("update", fixcost+varcost, "-", query)
	return &queryLocal{qcLocal: qcLocal{q: q, cost: fixcost + varcost, mode: qry.UpdateMode}}
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

// qcLocal - common base for queryLocal and cursorLocal

type qcLocal struct {
	// Query is embedded so most methods are "inherited" directly
	q    qry.Query
	keys []string // cache
	cost qry.Cost
	mode qry.Mode
}

func (q *qcLocal) Keys() []string {
	if q.keys == nil {
		keys := q.q.Keys()
		list := make([]string, len(keys))
		for i, k := range keys {
			list[i] = str.Join(",", k)
		}
		q.keys = list
	}
	return q.keys
}

func (q *qcLocal) Strategy(formatted bool) string {
	var strategy string
	if formatted {
		strategy = qry.Strategy(q.q) + "\n"
	} else {
		strategy = qry.String(q.q) + " "
	}
	n, _ := q.q.Nrows()
	return fmt.Sprint(strategy,
		"[nrecs~ ", n, " cost~ ", q.cost, " ", q.mode, "]")
}

func (q *qcLocal) Get(th *Thread, dir Dir) (Row, string) {
	defer th.Suneido.Store(th.Suneido.Swap(nil)) // use main Suneido object
	row := q.q.Get(th, dir)
	if row == nil {
		q.q.Rewind() // required for SuQuery to stick at eof unidirectionally
	}
	return row, q.q.Updateable()
}

func (q *qcLocal) Tree() Value {
	qry.CalcSelf(q.q)
	return qry.NewSuQueryNode(q.q)
}

func (q *qcLocal) Close() {
}

// cursorLocal

type cursorLocal struct {
	qcLocal
}

func (q *cursorLocal) Get(th *Thread, t ITran, dir Dir) (Row, string) {
	q.q.SetTran(t.(qry.QueryTran))
	return q.qcLocal.Get(th, dir)
}

func (q *cursorLocal) Rewind() {
	q.q.Rewind()
}

func (q *cursorLocal) Header() *Header {
	return q.q.Header()
}

func (q *cursorLocal) Order() []string {
	return q.q.Order()
}

// queryLocal

type queryLocal struct {
	qcLocal
	ra readAhead
}

func (q *queryLocal) Rewind() {
	if q.ra.pending.Load() {
		q.ra.mutex.Lock()
		defer q.ra.mutex.Unlock()
	}
	q.ra.fetch() // need to consume read ahead
	q.q.Rewind()
}

func (q *queryLocal) Header() *Header {
	if q.ra.pending.Load() {
		q.ra.mutex.Lock()
		defer q.ra.mutex.Unlock()
	}
	return q.q.Header().Dup()
}
func (q *queryLocal) Keys() []string {
	if q.ra.pending.Load() {
		q.ra.mutex.Lock()
		defer q.ra.mutex.Unlock()
	}
	return q.qcLocal.Keys()
}

func (q *queryLocal) Order() []string {
	if q.ra.pending.Load() {
		q.ra.mutex.Lock()
		defer q.ra.mutex.Unlock()
	}
	return q.q.Order()
}
func (q *queryLocal) Strategy(formatted bool) string {
	if q.ra.pending.Load() {
		q.ra.mutex.Lock()
		defer q.ra.mutex.Unlock()
	}
	return q.qcLocal.Strategy(formatted)
}

func (q *queryLocal) Output(th *Thread, rec Record) {
	q.q.Output(th, rec)
}

func (q *queryLocal) Tree() Value {
	if q.ra.pending.Load() {
		q.ra.mutex.Lock()
		defer q.ra.mutex.Unlock()
	}
	return q.qcLocal.Tree()
}

var readAheads atomic.Int32
var _ = AddInfo("database.readAheads", &readAheads)

func (q *queryLocal) Get(th *Thread, dir Dir) (Row, string) {
	ra := &q.ra
	ra.mutex.Lock()
	unlock := true
	defer func() {
		if unlock {
			ra.mutex.Unlock()
		}
	}()
	row, tbl, gotRow := ra.fetch()
	if gotRow && dir == Prev { // switched direction
		ra.disable = true  // prevent further read-aheads
		q.qcLocal.Get(th, Prev) // undo the read-ahead
		gotRow = false
	}
	if gotRow {
		readAheads.Add(1)
	} else {
		row, tbl = q.qcLocal.Get(th, dir)
	}
	if dir == Next && row != nil && !ra.disable && q.mode == qry.ReadMode {
		ra.pending.Store(true)
		unlock = false
		go func() {
			defer ra.mutex.Unlock()
			th := getPoolThread()
			defer threadPool.Put(th)
			ra.row, ra.tbl = q.qcLocal.Get(th, Next) // the actual read-ahead
		}()
	}
	return row, tbl
}

var threadPool sync.Pool

func getPoolThread() *Thread {
	if t := threadPool.Get(); t != nil {
		th := t.(*Thread)
		*th = Thread{} // clear the thread
		return th
	}
	return &Thread{}
}

type readAhead struct {
	pending atomic.Bool // true if read-ahead was triggered
	disable bool        // set if dir reverses to prevent further read-ahead
	mutex   sync.Mutex
	row     Row // the result of the read-ahead
	tbl     string
}

// fetch WARNING does not lock
func (ra *readAhead) fetch() (row Row, tbl string, gotRow bool) {
	if !ra.pending.Load() {
		return
	}
	row, tbl = ra.row, ra.tbl
	ra.row, ra.tbl = nil, ""
	ra.pending.Store(false)
	return row, tbl, true
}
