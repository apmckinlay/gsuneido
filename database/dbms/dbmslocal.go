package dbms

import (
	"fmt"
	"hash/adler32"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
)

// DbmsLocal implements the Dbms interface using a local database
// i.e. standalone
type DbmsLocal struct {
	libraries []string //TODO concurrency
}

func NewDbmsLocal() IDbms {
	return &DbmsLocal{}
}

// Dbms interface

var _ IDbms = (*DbmsLocal)(nil)

func (DbmsLocal) Admin(string) {
	panic("DbmsLocal Admin not implemented")
}

func (DbmsLocal) Auth(string) bool {
	panic("Auth only allowed on clients")
}

func (DbmsLocal) Check() string {
	panic("DbmsLocal Check not implemented")
}

func (DbmsLocal) Connections() Value {
	return EmptyObject
}

func (DbmsLocal) Cursors() int {
	panic("DbmsLocal Cursors not implemented")
}

func (DbmsLocal) Dump(string) string {
	panic("DbmsLocal Dump not implemented")
}

func (DbmsLocal) Exec(t *Thread, v Value) Value {
	fname := IfStr(ToObject(v).ListGet(0))
	if i := strings.IndexByte(fname, '.'); i != -1 {
		ob := Global.GetName(t, fname[:i])
		m := fname[i+1:]
		return t.CallMethodWithArgSpec(ob, m, ArgSpecEach1, v)
	}
	fn := Global.GetName(t, fname)
	return t.CallWithArgSpec(fn, ArgSpecEach1, v)
}

func (DbmsLocal) Final() int {
	panic("DbmsLocal Final not implemented")
}

func (DbmsLocal) Info() Value {
	panic("DbmsLocal Info not implemented")
}

func (DbmsLocal) Kill(string) int {
	panic("DbmsLocal Kill not implemented")
}

func (DbmsLocal) Load(string) int {
	panic("DbmsLocal Load not implemented")
}

func (DbmsLocal) LibGet(name string) (result []string) {
	// Temporary version that reads from text files
	defer func() {
		if e := recover(); e != nil {
			panic("error loading " + name + " " + fmt.Sprint(e))
			result = nil
		}
	}()
	dir := "../stdlib/"
	hash := adler32.Checksum([]byte(name))
	name = strings.ReplaceAll(name, "?", "Q")
	file := dir + name + "_" + strconv.FormatUint(uint64(hash), 16)
	s, err := ioutil.ReadFile(file)
	if err != nil {
		if !strings.HasPrefix(name, "Rule_") {
			fmt.Println("LOAD", file, "NOT FOUND")
		}
		return nil
	}
	// fmt.Println("LOAD", name, "SUCCEEDED")
	return []string{"stdlib", string(s)}
}

func (DbmsLocal) Libraries() *SuObject {
	return NewSuObject()
}

func (DbmsLocal) Log(s string) {
	f, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("can't open error.log " + err.Error())
	}
	if _, err := f.Write([]byte(s + "\n")); err != nil {
		panic("can't write to error.log " + err.Error())
	}
	f.Close()
}

func (DbmsLocal) Nonce() string {
	panic("nonce only allowed on clients")
}

func (DbmsLocal) Run(string) Value {
	panic("DbmsLocal Run not implemented")
}

var sessionId string

func (DbmsLocal) SessionId(id string) string {
	if id != "" {
		sessionId = id
	}
	return sessionId
}

func (DbmsLocal) Size() int64 {
	panic("DbmsLocal Size not implemented")
}

func (DbmsLocal) Token() string {
	panic("DbmsLocal Token not implemented")
}

var prevTimestamp SuDate

func (DbmsLocal) Timestamp() SuDate {
	t := Now()
	if t.Equal(prevTimestamp) {
		t = t.Plus(0, 0, 0, 0, 0, 0, 1)
	}
	prevTimestamp = t
	return t
}

func (DbmsLocal) Transactions() *SuObject {
	panic("DbmsLocal Transactions not implemented")
}

func (dl DbmsLocal) Unuse(lib string) bool {
	if lib == "stdlib" || !str.ListHas(dl.libraries, lib) {
		return false
	}
	dl.libraries = str.ListRemove(dl.libraries, lib)
	return true
}

func (dl DbmsLocal) Use(lib string) bool {
	if str.ListHas(dl.libraries, lib) {
		return false
	}
	//TODO check schema
	dl.libraries = append(dl.libraries, lib)
	return true
}

func (DbmsLocal) Close() {
}
