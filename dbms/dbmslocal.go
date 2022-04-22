// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/tools"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

// DbmsLocal implements the Dbms interface using a local database
// i.e. standalone
type DbmsLocal struct {
	db        *db19.Database
	libraries atomic.Value // []string
	badlibs   int32        // limits logging
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
	if sv == nil {
		sv = &dbms.db.Sviews
	}
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
	return EmptyObject
}

func (dbms *DbmsLocal) Cursor(query string, sv *Sviews) ICursor {
	if sv == nil {
		sv = &dbms.db.Sviews
	}
	q := qry.ParseQuery(query, dbms.db.NewReadTran(), sv)
	q, cost := qry.Setup(q, qry.CursorMode, dbms.db.NewReadTran())
	return cursorLocal{queryLocal{Query: q, cost: cost, mode: qry.CursorMode}}
}

func (*DbmsLocal) Cursors() int {
	return 0
}

func (dbms *DbmsLocal) DisableTrigger(table string) {
	dbms.db.DisableTrigger(table)
}
func (dbms *DbmsLocal) EnableTrigger(table string) {
	dbms.db.EnableTrigger(table)
}

func (dbms *DbmsLocal) Dump(table string) string {
	var err error
	if table == "" {
		_, _, err = tools.Dump(dbms.db, "database.su")
	} else {
		_, err = tools.DumpDbTable(dbms.db, table, table+".su")
	}
	if err != nil {
		return err.Error()
	}
	return ""
}

func (*DbmsLocal) Exec(t *Thread, v Value) Value {
	return t.RunWithMainSuneido(func() Value {
		trace.Dbms.Println("Exec", v)
		fname := ToStr(ToContainer(v).ListGet(0))
		if i := strings.IndexByte(fname, '.'); i != -1 {
			ob := Global.GetName(t, fname[:i])
			m := fname[i+1:]
			return t.CallLookupEach1(ob, m, v)
		}
		fn := Global.GetName(t, fname)
		return t.CallEach1(fn, v)
	})
}

func (*DbmsLocal) Final() int {
	return 0 //TODO
}

// Get implements QueryFirst, QueryLast, Query1
func (dbms *DbmsLocal) Get(
	th *Thread, query string, dir Dir, sv *Sviews) (Row, *Header, string) {
	tran := dbms.db.NewReadTran()
	defer tran.Complete()
	if sv == nil {
		sv = &dbms.db.Sviews
	}
	return get(th, tran, query, dir, sv)
}

func get(th *Thread, tran qry.QueryTran, query string, dir Dir,
	sv *Sviews) (Row, *Header, string) {
	q := qry.ParseQuery(query, tran, sv)
	q, _ = qry.Setup(q, qry.ReadMode, tran)
	only := false
	if dir == Only {
		only = true
		dir = Next
	}
	row := q.Get(th, dir)
	if row == nil {
		return nil, nil, ""
	}
	if only && q.Get(th, dir) != nil {
		panic("Query1 not unique: " + query)
	}
	return row, q.Header(), q.Updateable()
}

func (dbms *DbmsLocal) Info() Value {
	ob := &SuObject{}
	ob.Set(SuStr("currentSize"), Int64Val(int64(dbms.db.Size())))
	return ob
}

func (*DbmsLocal) Kill(addr string) int {
	if options.Action == "server" {
		return kill(addr)
	}
	return 0
}

func (dbms *DbmsLocal) Load(table string) int {
	n, err := tools.LoadDbTable(table, dbms.db)
	if err != nil {
		panic(err.Error())
	}
	return n
}

func (dbms *DbmsLocal) LibGet(name string) []string {
	defer func() {
		if e := recover(); e != nil {
			// debug.PrintStack()
			panic("error loading " + name + " " + fmt.Sprint(e))
		}
	}()

	results := make([]string, 0, 2)
	rt := dbms.db.NewReadTran()
	libs := dbms.libraries.Load().([]string)
	for _, lib := range libs {
		s := dbms.libGet(rt, lib, name)
		if s != "" {
			results = append(results, lib, string(s))
		}
	}
	return results
}

var libKey = []string{"name", "group"} // const

func (dbms *DbmsLocal) libGet(rt *db19.ReadTran, lib, name string) string {
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
	if atomic.AddInt32(&dbms.badlibs, 1) <= 1 {
		log.Println("ERROR: invalid library: " + lib)
	}
}

func (dbms *DbmsLocal) Libraries() []string {
	// library list is not mutated so it's thread safe to return
	return dbms.libraries.Load().([]string)
}

func (*DbmsLocal) Log(s string) {
	log.Println(s)
}

func (*DbmsLocal) Nonce(th *Thread) string {
	th.Nonce = Nonce()
	return th.Nonce
}

func (*DbmsLocal) Run(th *Thread, s string) Value {
	return th.RunWithMainSuneido(func() Value {
		trace.Dbms.Println("Run", s)
		return compile.EvalString(th, s)
	})
}

func (dbms *DbmsLocal) Schema(table string) string {
	return dbms.db.Schema(table)
}

func (*DbmsLocal) SessionId(t *Thread, id string) string {
	if id != "" {
		t.SetSession(id)
	}
	return t.Session()
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
		return nil
	}
	return &ReadTranLocal{ReadTran: dbms.db.NewReadTran()}
}

// Transctions only returns the update transactions
func (dbms *DbmsLocal) Transactions() *SuObject {
	var ob SuObject
	trans := dbms.db.Transactions()
	for _, t := range trans {
		ob.Add(IntVal(t<<1 | 1)) // update tran# are odd
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
	oldlibs := dbms.libraries.Load().([]string)
	newlibs := fn(oldlibs)
	if newlibs == nil {
		return false
	}
	atomic.StoreInt32(&dbms.badlibs, 0) // reset logging
	return slices.Equal(oldlibs, dbms.libraries.Swap(newlibs).([]string))
}

func (dbms *DbmsLocal) Unwrap() IDbms {
	return dbms
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

type TranLocal interface {
	Num() int
}

type ReadTranLocal struct {
	*db19.ReadTran
	TranLocal
}

func (t ReadTranLocal) Get(th *Thread, query string, dir Dir,
	sv *Sviews) (Row, *Header, string) {
	if sv == nil {
		sv = t.GetSviews()
	}
	return get(th, t.ReadTran, query, dir, sv)
}

func (t ReadTranLocal) Query(query string, sv *Sviews) IQuery {
	if sv == nil {
		sv = t.GetSviews()
	}
	q := qry.ParseQuery(query, t.ReadTran, sv)
	q, cost := qry.Setup(q, qry.ReadMode, t.ReadTran)
	return queryLocal{Query: q, cost: cost, mode: qry.ReadMode}
}

func (t ReadTranLocal) Action(*Thread, string, *Sviews) int {
	panic("cannot do action in read-only transaction")
}

// UpdateTranLocal --------------------------------------------------------

type UpdateTranLocal struct {
	*db19.UpdateTran
	TranLocal
}

func (t UpdateTranLocal) Get(th *Thread, query string, dir Dir,
	sv *Sviews) (Row, *Header, string) {
	if sv == nil {
		sv = t.GetSviews()
	}
	return get(th, t.UpdateTran, query, dir, sv)
}

func (t UpdateTranLocal) Query(query string, sv *Sviews) IQuery {
	if sv == nil {
		sv = t.GetSviews()
	}
	q := qry.ParseQuery(query, t.UpdateTran, sv)
	q, cost := qry.Setup(q, qry.UpdateMode, t.UpdateTran)
	return queryLocal{Query: q, cost: cost, mode: qry.UpdateMode}
}

func (t UpdateTranLocal) Action(th *Thread, action string, sv *Sviews) int {
	trace.Dbms.Println("Action", action)
	if sv == nil {
		sv = t.GetSviews()
	}
	return qry.DoAction(th, t.UpdateTran, action, sv)
}

// queryLocal

type queryLocal struct {
	// Query is embedded so most methods are "inherited" directly
	qry.Query
	cost qry.Cost
	mode qry.Mode
	keys []string // cache
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

func (q queryLocal) Strategy() string {
	return fmt.Sprint(q.String(),
		" [nrecs~ ", q.Nrows(), " cost~ ", q.cost, " ", q.mode, "]")
}

func (q queryLocal) Order() []string {
	return q.Query.Ordering()
}

func (q queryLocal) Get(th *Thread, dir Dir) (Row, string) {
	row := q.Query.Get(th, dir)
	if row == nil {
		// this is required for SuQuery to stick at eof unidirectionally
		q.Query.Rewind()
	}
	return row, q.Query.Updateable()
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
