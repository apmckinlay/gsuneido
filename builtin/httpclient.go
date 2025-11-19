// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/dbms"
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
	cl := addHeader(req.Header, args[3])
	if cl != -1 {
		req.ContentLength = cl
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
		return &suHttpResponse{resp: resp, body: hacks.BStoS(body)}
	}
}

var _ = builtin(HttpsClient, `(method, url, 
	content = '', header = #(), timeout = 60, block = false)`)

// HttpsClient makes an HTTPS request embedded cert
func HttpsClient(th *Thread, args []Value) Value {
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(dbms.ServerCert)
	if !ok {
		panic("Failed to append embedded cert to pool")
	}
	config := &tls.Config{
		RootCAs:    caCertPool,
		ServerName: "localhost", // Must match CN or SAN
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}
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
	cl := addHeader(req.Header, args[3])
	if cl != -1 {
		req.ContentLength = cl
	}

	resp, err := client.Do(req)
	hck(err)
	defer resp.Body.Close()
	if block := args[5]; block != False {
		th.Call(block, &suHttpResponse{resp: resp})
		return nil
	} else {
		body, err := io.ReadAll(resp.Body)
		hck(err)
		return &suHttpResponse{resp: resp, body: hacks.BStoS(body)}
	}
}

func addHeader(header http.Header, h Value) (contentLength int64) {
	contentLength = -1
	if hdr, ok := h.ToContainer(); ok {
		it := hdr.Iter2(false, true)
		for k, v := it(); k != nil; k, v = it() {
			key := strings.ReplaceAll(ToStr(k), "_", "-")
			if str.EqualCI(key, "Content-Length") {
				contentLength = int64(ToInt(v))
			}
			if vals, ok := v.ToContainer(); ok {
				list := make([]string, vals.ListSize())
				for i := 0; i < vals.ListSize(); i++ {
					list[i] = AsStr(vals.ListGet(i))
				}
				header[key] = list
			} else {
				header[key] = []string{AsStr(v)}
			}
		}
	}
	return
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
	body string
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

func (hr *suHttpResponse) Get(_ *Thread, k Value) Value {
	key := ToStr(k)
	switch key {
	case "status":
		return IntVal(hr.resp.StatusCode)
	case "proto":
		return SuStr(hr.resp.Proto)
	case "header": // DEPRECATED
		var sb strings.Builder
		fmt.Fprintln(&sb, hr.resp.Proto, hr.resp.Status)
		hr.resp.Header.Write(&sb)
		return SuStr(sb.String())
	case "content", "body":
		return SuStr(hr.body)
	default:
		key = strings.ReplaceAll(key, "_", "-")
		v := hr.resp.Header[key]
		if len(v) == 1 {
			return SuStr(v[0])
		} else if len(v) > 1 {
			return SuObjectOfStrs(v)
		}
	}
	return EmptyStr
}

var suHttpResponseMethods = methods("resp")

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
