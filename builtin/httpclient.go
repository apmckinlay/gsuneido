// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
)

var _ = builtin(HttpClient2, `(method, url, 
	content = '', header = #(), timeout = 60, block = false)`)

// HttpClient2 is a wrapper around Go net/http
func HttpClient2(th *Thread, args []Value) Value {
	method := ToStr(args[0])
	url := ToStr(args[1])
	var rdr io.Reader
	if isFunction(args[2]) {
		rdr = &reader{th: th, fn: args[2]}
	} else if content := ToStr(args[2]); content != "" {
		rdr = strings.NewReader(content)
	}

	to := time.Duration(ToInt(args[4])) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), to)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, rdr)
	hck(err)
	req.Header.Set("User-Agent", "Suneido")
	hdr := ToContainer(args[3])
	f := hdr.Iter2(false, true)
	for k, v := f(); k != nil; k, v = f() {
		key := strings.ReplaceAll(ToStr(k), "_", "-")
		if str.EqualCI(key, "Content-Length") {
			req.ContentLength = int64(ToInt(v))
		} else {
			req.Header.Set(key, AsStr(v))
		}
	}

	resp, err := http.DefaultClient.Do(req)
	hck(err)
	defer resp.Body.Close()
	if block := args[5]; block != False {
		th.Call(block, &suHttpResponse{resp: resp})
		return nil
	} else {
		body, err := io.ReadAll(resp.Body)
		hck(err)
		ob := &SuObject{}
		ob.Set(SuStr("content"), SuStr(hacks.BStoS(body)))
		ob.Set(SuStr("header"), SuStr(headerString(resp)))
		return ob
	}
}

func headerString(resp *http.Response) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, resp.Proto, resp.Status)
	resp.Header.Write(&sb)
	return sb.String()
}

func hck(err error) {
	if err != nil {
		panic("HttpClient: " + err.Error())
	}
}

// reader adapts a Suneido callable to io.Reader
type reader struct {
	th *Thread
	fn Value
}

func (r *reader) Read(buf []byte) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("HttpClient: %v", r)
		}
	}()
	v := r.th.Call(r.fn, IntVal(len(buf)))
	s := ToStr(v)
	n = copy(buf, s)
	if n == 0 {
		err = io.EOF
	}
	return
}

//-------------------------------------------------------------------

// suHttpResponse wraps an http.Response for Suneido to access
// @immutable
type suHttpResponse struct {
	ValueBase[suHttpResponse]
	resp *http.Response
}

var _ Value = (*suHttpResponse)(nil)

func (hr *suHttpResponse) Equal(other any) bool {
	return hr == other
}

func (hr *suHttpResponse) Lookup(_ *Thread, method string) Value {
	return suHttpResponseMethods[method]
}

func (hr *suHttpResponse) SetConcurrent() {
	// safe
}

var suHttpResponseMethods = methods("resp")

var _ = method(resp_Header, "()")

func resp_Header(this Value) Value {
	return SuStr(headerString(this.(*suHttpResponse).resp))
}

var _ = method(resp_Status, "()")

func resp_Status(this Value) Value {
	return IntVal(this.(*suHttpResponse).resp.StatusCode)
}

var _ = method(resp_Read, "(n)")

func resp_Read(this Value, a Value) Value {
	rdr := this.(*suHttpResponse).resp.Body
	n := ToInt(a)
	if n > readMax {
		panic("HttpClient: Read too large")
	}
	buf := make([]byte, n)
	nr, err := rdr.Read(buf)
	if nr > 0 {
		return SuStr(hacks.BStoS(buf[:nr]))
	}
	if err != nil {
		if err == io.EOF {
			return False
		}
		panic(fmt.Sprint("HttpClient: ", err))
	}
	return EmptyStr
}

var _ = method(resp_CopyTo, "(dest, nbytes = false)")

func resp_CopyTo(th *Thread, this Value, args []Value) Value {
	rd := this.(*suHttpResponse).resp.Body
	return CopyTo(th, rd, args[0], args[1])
}
