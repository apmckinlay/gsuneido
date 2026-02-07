// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !gui

package builtin

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/dbms"
)

var _ = builtin(HttpsServer, `(port, app, stop = false)`)

func HttpsServer(th *Thread, args []Value) Value {
	guardSandbox("HttpsServer")
	cert, err := tls.X509KeyPair(dbms.ServerCert, dbms.ServerKey)
	if err != nil {
		log.Fatalf("ERROR: Failed to load embedded key pair: %v", err)
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	port := ToInt(args[0])
	addr := fmt.Sprint(":", port)
	server := &http.Server{
		Addr:      addr,
		Handler:   &HttpHandler{th: th, app: args[1]},
		TLSConfig: tlsConfig,
	}
	if ob, ok := args[2].ToContainer(); ok {
		ob.Put(th, SuStr("stop"), &suStopper{server: server})
	}
	if err := server.ListenAndServeTLS("", ""); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprint("HttpServer:", err))
	}
	return nil
}
