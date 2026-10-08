// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !gui

package builtin

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/dbms"
)

var _ = builtin(HttpsServer, `(port, app, stop = false) :void`)

func HttpsServer(th *Thread, args []Value) Value {
	guardSandbox("HttpsServer")
	cert, err := tls.X509KeyPair(dbms.ServerCert, dbms.ServerKey)
	if err != nil {
		log.Fatalf("ERROR: Failed to load embedded key pair: %v", err)
	}
	clientCAPool := x509.NewCertPool()
	if !clientCAPool.AppendCertsFromPEM(dbms.ClientCert) {
		log.Fatalf("ERROR: Failed to append embedded client cert to pool")
	}
	// mutual TLS: only accept clients that present the embedded client cert
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAPool,
	}
	port := ToInt(args[0])
	addr := fmt.Sprint(":", port)
	server := &http.Server{
		Addr:      addr,
		Handler:   newHttpHandler(th, args[1]),
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
