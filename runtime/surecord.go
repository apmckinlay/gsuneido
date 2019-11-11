package runtime

import (
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/pack"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/verify"
)

// SuRecord is an SuObject with observers and rules and a default value of ""
type SuRecord struct {
	ob SuObject
	CantConvert
	// observers is from record.Observer(fn)
	observers ValueList
	// invalidated accumulates keys needing observers called
	invalidated str.Queue
	// invalid is the fields that need to be recalculated
	invalid map[string]bool
	// dependents are the fields that depend on a field
	dependents map[string][]string
	// activeObservers is used to prevent infinite recursion
	activeObservers ActiveObserverList
	// attachedRules is from record.AttachRule(key,fn)
	attachedRules map[string]Value

	// row is used when it is from the database
	row Row
	// header is the Header for row
	hdr *Header
	// tran is the database transaction used to read the record
	tran *SuTran
	// recadr is the record address in the database
	recadr int
	// unpacked is true if the row has been unpacked into ob
	unpacked bool
}

//go:generate genny -in ../../GoTemplates/list/list.go -out alist.go -pkg runtime gen "V=activeObserver"
//go:generate genny -in ../../GoTemplates/list/list.go -out vlist.go -pkg runtime gen "V=Value"

func NewSuRecord() *SuRecord {
	return &SuRecord{ob: SuObject{defval: EmptyStr}}
}

// SuRecordFromObject creates a record from an arguments object
// WARNING: it does not copy the data, the original object should be discarded
func SuRecordFromObject(ob *SuObject) *SuRecord {
	return &SuRecord{
		ob: SuObject{list: ob.list, named: ob.named, defval: EmptyStr}}
}

func SuRecordFromRow(row Row, hdr *Header, tran *SuTran) *SuRecord {
	if hdr.Map == nil { //TODO concurrency
		hdr.Map = make(map[string]RowAt, len(hdr.Fields))
		for ri, r := range hdr.Fields {
			for fi, f := range r {
				hdr.Map[f] = RowAt{int16(ri), int16(fi)}
			}
		}
	}
	dependents := deps(row, hdr)
	return &SuRecord{row: row, hdr: hdr, tran: tran, recadr: row[0].Adr,
		ob: SuObject{defval: EmptyStr}, dependents: dependents}
}

func deps(row Row, hdr *Header) map[string][]string {
	dependents := map[string][]string{}
	for _, f := range hdr.Fields[0] {
		if strings.HasSuffix(f, "_deps") {
			deps := strings.Split(ToStr(row.Get(hdr, f)), ",")
			f = f[:len(f)-5]
			for _, d := range deps {
				if !str.ListHas(dependents[d], f) {
					dependents[d] = append(dependents[d], f)
				}
			}
		}
	}
	return dependents
}

func (r *SuRecord) Copy() Container {
	return r.slice(0)
}

func (r *SuRecord) slice(n int) *SuRecord {
	// keep row and hdr even if unpacked, to help ToRecord
	return &SuRecord{
		ob:         r.ob.slice(n),
		row:        r.row,
		hdr:        r.hdr,
		unpacked:	r.unpacked,
		dependents: r.copyDeps(),
		invalid:    r.copyInvalid()}
}

func (r *SuRecord) copyDeps() map[string][]string {
	copy := make(map[string][]string, len(r.dependents))
	for k, v := range r.dependents {
		copy[k] = append(v[:0:0], v...) // copy slice
	}
	return copy
}

func (r *SuRecord) copyInvalid() map[string]bool {
	copy := make(map[string]bool, len(r.invalid))
	for k, v := range r.invalid {
		copy[k] = v
	}
	return copy
}

func (*SuRecord) Type() types.Type {
	return types.Record
}

func (r *SuRecord) String() string {
	s := r.ToObject().String()
	return "[" + s[2:len(s)-1] + "]"
}

func (r *SuRecord) Show() string {
	s := r.ob.Show()
	return "[" + s[2:len(s)-1] + "]"
}

func (*SuRecord) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Record")
}

func (r *SuRecord) Compare(other Value) int {
	return r.ToObject().Compare(other)
}

func (r *SuRecord) Equal(other interface{}) bool {
	return r.ToObject().Equal(other)
}

func (r *SuRecord) Hash() uint32 {
	return r.ob.Hash()
}

func (r *SuRecord) Hash2() uint32 {
	return r.ob.Hash2()
}

func (r *SuRecord) RangeTo(from int, to int) Value {
	return r.ob.RangeTo(from, to)
}

func (r *SuRecord) RangeLen(from int, n int) Value {
	return r.ob.RangeLen(from, n)
}

func (r *SuRecord) ToContainer() (Container, bool) {
	return r, true
}

// Container --------------------------------------------------------

var _ Container = (*SuRecord)(nil)

func (r *SuRecord) ToObject() *SuObject {
	if r.row != nil && !r.unpacked {
		for ri, rf := range r.hdr.Fields {
			for fi, f := range rf {
				if f != "-" && !strings.HasSuffix(f, "_deps") {
					key := SuStr(f)
					if !r.ob.HasKey(key) {
						if val := r.row[ri].GetRaw(fi); val != "" {
							r.ob.Set(key, Unpack(val))
						}
					}
				}
			}
		}
		r.unpacked = true
		// keep row and hdr for ToRecord
	}
	return &r.ob
}

func (r *SuRecord) Add(val Value) {
	r.ob.Add(val)
}

func (r *SuRecord) Insert(at int, val Value) {
	r.ob.Insert(at, val)
}

func (r *SuRecord) HasKey(key Value) bool {
	if r.ob.HasKey(key) {
		return true
	}
	if r.row != nil {
		if k, ok := key.ToStr(); ok {
			return r.row.GetRaw(r.hdr, k) != ""
		}
	}
	return false
}

func (r *SuRecord) Set(key Value, val Value) {
	r.Put(nil, key, val)
}

func (r *SuRecord) Clear() {
	r.ob.mustBeMutable()
	*r = SuRecord{}
}

func (r *SuRecord) SetReadOnly() {
	r.ob.SetReadOnly()
}

func (r *SuRecord) IsReadOnly() bool {
	return r.ob.IsReadOnly()
}

func (r *SuRecord) IsNew() bool {
	return r.row == nil
}

func (r *SuRecord) Delete(t *Thread, key Value) bool {
	r.ob.mustBeMutable()
	r.ToObject()
	if r.ob.Delete(t, key) {
		if keystr, ok := key.ToStr(); ok {
			r.invalidateDependents(keystr)
			r.callObservers(t, keystr)
		}
		return true
	}
	return false
}

func (r *SuRecord) Erase(t *Thread, key Value) bool {
	r.ob.mustBeMutable()
	r.ToObject()
	if r.ob.Erase(t, key) {
		if keystr, ok := key.ToStr(); ok {
			r.invalidateDependents(keystr)
			r.callObservers(t, keystr)
		}
		return true
	}
	return false
}

func (r *SuRecord) ListSize() int {
	return r.ob.ListSize()
}

func (r *SuRecord) ListGet(i int) Value {
	return r.ob.ListGet(i)
}

func (r *SuRecord) NamedSize() int {
	if r.row != nil && !r.unpacked {
		return r.rowSize()
	}
	return r.ob.NamedSize()
}

func (r *SuRecord) rowSize() int {
	n := r.ob.NamedSize()
	for ri, rf := range r.hdr.Fields {
		for fi, f := range rf {
			if f != "-" && !strings.HasSuffix(f, "_deps") {
				key := SuStr(f)
				if !r.ob.HasKey(key) {
					if val := r.row[ri].GetRaw(fi); val != "" {
						n++
					}
				}
			}
		}
	}
	return n
}

func (r *SuRecord) ArgsIter() func() (Value, Value) {
	return r.ToObject().ArgsIter()
}

func (r *SuRecord) Iter2(list bool, named bool) func() (Value, Value) {
	return r.ToObject().Iter2(list, named)
}

func (r *SuRecord) Slice(n int) Container {
	return r.slice(n)
}

func (r *SuRecord) Iter() Iter {
	return &obIter{ob: &r.ob, iter: r.ob.Iter2(true, true),
		result: func(k, v Value) Value { return v }}
}

// ------------------------------------------------------------------

func (r *SuRecord) Put(t *Thread, keyval Value, val Value) {
	if key, ok := keyval.ToStr(); ok {
		delete(r.invalid, key)
		old := r.ob.GetIfPresent(t, keyval)
		r.ob.Set(keyval, val)
		if old != nil && val.Equal(old) {
			return
		}
		r.invalidateDependents(key)
		r.callObservers(t, key)
	} else {
		r.ob.Set(keyval, val)
	}
}

func (r *SuRecord) invalidateDependents(key string) {
	for _, d := range r.dependents[key] {
		r.invalidate(d)
	}
}

func (r *SuRecord) Invalidate(t *Thread, key string) {
	r.invalidate(key)
	r.callObservers(t, key)
}

func (r *SuRecord) invalidate(key string) {
	if r.invalid[key] {
		return
	}
	r.invalidated.Add(key) // for observers
	if r.invalid == nil {
		r.invalid = make(map[string]bool)
	}
	r.invalid[key] = true
	r.invalidateDependents(key)
}

func (r *SuRecord) PreSet(key, val Value) {
	r.ob.Set(key, val)
}

func (r *SuRecord) Observer(ofn Value) {
	r.observers.Push(ofn)
}

func (r *SuRecord) RemoveObserver(ofn Value) bool {
	return r.observers.Remove(ofn)
}

func (r *SuRecord) callObservers(t *Thread, key string) {
	r.callObservers2(t, key)
	for !r.invalidated.Empty() {
		if k := r.invalidated.Take(); k != key {
			r.callObservers2(t, k)
		}
	}
}

func (r *SuRecord) callObservers2(t *Thread, key string) {
	for _, x := range r.observers.list {
		ofn := x.(Value)
		if !r.activeObservers.Has(activeObserver{ofn, key}) {
			func(ofn Value, key string) {
				r.activeObservers.Push(activeObserver{ofn, key})
				defer r.activeObservers.Pop()
				t.pushCall(ofn, r, argSpecMember, SuStr(key))
			}(ofn, key)
		}
	}
}

var argSpecMember = &ArgSpec{Nargs: 1,
	Spec: []byte{0}, Names: []Value{SuStr("member")}}

type activeObserver struct {
	obs Value
	key string
}

func (a activeObserver) Equal(other activeObserver) bool {
	return a.key == other.key && a.obs.Equal(other.obs)
}

// ------------------------------------------------------------------

// Get returns the value associated with a key, or defval if not found
func (r *SuRecord) Get(t *Thread, key Value) Value {
	if val := r.GetIfPresent(t, key); val != nil {
		return val
	}
	return r.ob.defaultValue(key)
}

// GetIfPresent is the same as Get
// except it returns nil instead of defval for missing members
func (r *SuRecord) GetIfPresent(t *Thread, keyval Value) Value {
	result := r.ob.GetIfPresent(t, keyval)
	if key, ok := keyval.ToStr(); ok {
		// only do record stuff when key is a string
		if result == nil && r.row != nil {
			raw := r.row.GetRaw(r.hdr, key)
			if raw != "" {
				val := Unpack(raw)
				r.PreSet(keyval, val) // cache unpacked value
				return val
			}
		}
		if t != nil {
			if ar := t.rules.top(); ar.rec == r { // identity (not Equal)
				r.addDependent(ar.key, key)
			}
		}
		if result == nil || r.invalid[key] {
			if x := r.getSpecial(key); x != nil {
				result = x
			} else if x = r.callRule(t, key); x != nil {
				result = x
			}
		}
	}
	return result
}

// GetPacked is used by ToRecord to build a Record for the database.
// It is like Get except it returns the value packed,
// using the already packed value from the row when possible.
// It does not add dependencies or handle special fields (e.g. _lower!)
func (r *SuRecord) GetPacked(t *Thread, key string) string {
	result := r.ob.GetIfPresent(t, SuStr(key))
	if result == nil && r.row != nil {
		if s := r.row.GetRaw(r.hdr, key); s != "" {
			return s
		}
	}
	if result == nil || r.invalid[key] {
		if x := r.callRule(t, key); x != nil {
			result = x
		}
	}
	if result == nil {
		result = r.ob.defval
		if result == nil {
			return ""
		}
	}
	return PackValue(result)
}

func (r *SuRecord) addDependent(from, to string) {
	if from == to {
		return
	}
	if r.dependents == nil {
		r.dependents = make(map[string][]string)
	}
	r.dependents[to] = append(r.dependents[to], from)
}

func (r *SuRecord) getSpecial(key string) Value {
	if strings.HasSuffix(key, "_lower!") {
		key = key[0 : len(key)-7]
		if val := r.ob.GetIfPresent(nil, SuStr(key)); val != nil {
			if vs, ok := val.ToStr(); ok {
				val = SuStr(strings.ToLower(vs))
			}
			return val
		}
	}
	return nil
}

func (r *SuRecord) callRule(t *Thread, key string) Value {
	delete(r.invalid, key)
	if rule := r.getRule(t, key); rule != nil && !t.rules.has(r, key) {
		val := r.catchRule(t, rule, key)
		if val != nil && !r.ob.IsReadOnly() {
			r.PreSet(SuStr(key), val)
		}
		return val
	}
	return nil
}

func (r *SuRecord) catchRule(t *Thread, rule Value, key string) Value {
	t.rules.push(r, key)
	defer func() {
		t.rules.pop()
		if e := recover(); e != nil {
			panic(toStr(e) + " (rule for " + key + ")")
		}
	}()
	return t.CallThis(rule, r)
}

// activeRules stack
type activeRules struct {
	list []activeRule
}
type activeRule struct {
	rec *SuRecord
	key string
}

func (ar *activeRules) push(r *SuRecord, key string) {
	ar.list = append(ar.list, activeRule{r, key})
}
func (ar *activeRules) top() activeRule {
	if len(ar.list) == 0 {
		return activeRule{}
	}
	return ar.list[len(ar.list)-1]
}
func (ar *activeRules) pop() {
	ar.list = ar.list[:len(ar.list)-1]
}
func (ar *activeRules) has(r *SuRecord, key string) bool {
	for _, x := range ar.list {
		if x.rec == r && x.key == key { // record identity (not Equal)
			return true
		}
	}
	return false
}

type strable interface{ String() string }
type errable interface{ Error() string }

func toStr(e interface{}) string {
	if s, ok := e.(string); ok {
		return s
	}
	if v, ok := e.(Value); ok {
		return AsStr(v)
	}
	if sa, ok := e.(strable); ok {
		return sa.String()
	}
	if ea, ok := e.(errable); ok {
		return ea.Error()
	}
	return "???"
}

func (r *SuRecord) getRule(t *Thread, key string) Value {
	if rule, ok := r.attachedRules[key]; ok {
		verify.That(rule != nil)
		return rule
	}
	if r.ob.defval != nil && t != nil {
		return Global.FindName(t, "Rule_"+key)
	}
	return nil
}

func (r *SuRecord) AttachRule(key, callable Value) {
	if r.attachedRules == nil {
		r.attachedRules = make(map[string]Value)
	}
	r.attachedRules[AsStr(key)] = callable
}

func (r *SuRecord) GetDeps(key string) Value {
	var sb strings.Builder
	sep := ""
	for to, froms := range r.dependents {
		for _, from := range froms {
			if from == key {
				sb.WriteString(sep)
				sb.WriteString(to)
				sep = ","
			}
		}
	}
	return SuStr(sb.String())
}

func (r *SuRecord) SetDeps(key, deps string) {
	if deps == "" {
		return
	}
outer:
	for _, to := range strings.Split(deps, ",") {
		to = strings.TrimSpace(to)
		for _, from := range r.dependents[to] {
			if from == key {
				continue outer
			}
		}
		r.addDependent(key, to)
	}
}

func (r *SuRecord) Transaction() *SuTran {
	return r.tran
}

// ToRecord converts this SuRecord to a Record to be stored in the database
func (r *SuRecord) ToRecord(t *Thread, hdr *Header) Record {
	fields := hdr.Fields[0]

	// access all the fields to ensure dependencies are created
	for _, f := range fields {
		//TODO optimize this - don't force unpack/pack on every field
		r.Get(t, SuStr(f))
	}

	// invert stored dependencies
	deps := map[string][]string{}
	for k, v := range r.dependents {
		for _, d := range v {
			d_deps := d + "_deps"
			if str.ListHas(fields, d_deps) {
				deps[d_deps] = append(deps[d_deps], k)
			}
		}
	}

	rb := RecordBuilder{}
	var tsField string
	var ts SuDate
	for _, f := range fields {
		if f == "" {
			rb.AddRaw("")
		} else if strings.HasSuffix(f, "_TS") { // also done in SuObject ToRecord
			tsField = f
			ts = t.Dbms().Timestamp()
			rb.Add(ts)
		} else if d, ok := deps[f]; ok {
			rb.Add(SuStr(strings.Join(d, ",")))
		} else {
			rb.AddRaw(r.GetPacked(t, f))
		}
	}
	if tsField != "" && !r.IsReadOnly() {
		r.Put(t, SuStr(tsField), ts)
	}
	return rb.Build()
}

// RecordMethods is initialized by the builtin package
var RecordMethods Methods

var gnRecords = Global.Num("Records")

func (*SuRecord) Lookup(t *Thread, method string) Callable {
	if m := Lookup(t, RecordMethods, gnRecords, method); m != nil {
		return m
	}
	return (*SuObject)(nil).Lookup(t, method)
}

// Packable ---------------------------------------------------------

var _ Packable = (*SuRecord)(nil)

func (r *SuRecord) PackSize(clock *int32) int {
	return r.ob.PackSize(clock)
}

func (r *SuRecord) PackSize2(clock int32, stack packStack) int {
	return r.ob.PackSize2(clock, stack)
}

func (r *SuRecord) PackSize3() int {
	return r.ob.PackSize3()
}

func (r *SuRecord) Pack(clock int32, buf *pack.Encoder) {
	r.ob.pack(clock, buf, PackRecord)
}

func UnpackRecord(s string) *SuRecord {
	r := NewSuRecord()
	unpackObject(s, &r.ob)
	return r
}

// old format

func UnpackRecordOld(s string) *SuRecord {
	r := NewSuRecord()
	unpackObjectOld(s, &r.ob)
	return r
}

// database

func (r *SuRecord) DbDelete() {
	r.ckModify("Delete")
	r.tran.Erase(r.recadr)
}

func (r *SuRecord) DbUpdate(t *Thread, ob Container) {
	r.ckModify("Update")
	var rec = ob.ToRecord(t, r.hdr)
	r.recadr = r.tran.Update(r.recadr, rec)
}

func (r *SuRecord) ckModify(op string) {
	if r.row == nil {
		panic("record." + op + ": not a database record")
	}
	if r.tran == nil {
		panic("record." + op + ": no Transaction")
	}
}
