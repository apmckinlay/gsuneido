package clientserver

import (
	"fmt"
	"hash/adler32"
	"io/ioutil"
	"strconv"
	"strings"
)

// DbmsLocal implements the Dbms interface using a local database
// i.e. standalone
type DbmsLocal struct {
}

func NewDbmsLocal() Dbms {
	return &DbmsLocal{}
}

// Dbms interface

var _ Dbms = (*DbmsLocal)(nil)

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
