// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/system"

	//lint:ignore SA1019 best option
	"golang.org/x/crypto/openpgp"
	"golang.org/x/exp/maps"
)

type suOpenPGP struct {
	staticClass[suOpenPGP]
}

func init() {
	Global.Builtin("OpenPGP", &suOpenPGP{})
}

func (*suOpenPGP) String() string {
	return "OpenPGP /* builtin class */"
}

func (*suOpenPGP) Lookup(_ *Thread, method string) Callable {
	return openpgpMethods[method]
}

var openpgpMethods = methods()

var _ = staticMethod(opgp_Members, "()")

func opgp_Members() Value {
	return SuObjectOfStrs(maps.Keys(openpgpMethods))
}

var _ = staticMethod(opgp_SymmetricEncrypt,
	"(passphrase, source, toFile = false)")

func opgp_SymmetricEncrypt(passphrase, source, toFile Value) Value {
	if toFile == False {
		return symStr(passphrase, source, symEncrypt)
	}
	return symFile(passphrase, source, toFile, symEncrypt)
}

var _ = staticMethod(opgp_SymmetricDecrypt,
	"(passphrase, source, toFile = false)")

func opgp_SymmetricDecrypt(passphrase, source, toFile Value) Value {
	if toFile == False {
		return symStr(passphrase, source, symDecrypt)
	}
	return symFile(passphrase, source, toFile, symDecrypt)
}

type encdec func(passphrase string, src io.Reader, dst io.Writer)

func symStr(passphrase Value, source Value, f encdec) Value {
	src := strings.NewReader(ToStr(source))
	dst := new(bytes.Buffer)
	f(ToStr(passphrase), src, dst)
	return SuStr(dst.String())
}

func symFile(passphrase Value, fromFile, toFile Value, f encdec) Value {
	src, err := os.Open(ToStr(fromFile))
	ck(err)
	defer src.Close()
	to := ToStr(toFile)
	dst, err := os.CreateTemp(filepath.Dir(to), "su")
	ck(err)
	defer os.Remove(dst.Name())
	defer dst.Close()
	f(ToStr(passphrase), src, dst)
	dst.Close()
	system.RenameBak(dst.Name(), to)
	return nil
}

func symEncrypt(passphrase string, src io.Reader, dst io.Writer) {
	encrypter, err := openpgp.SymmetricallyEncrypt(dst, []byte(passphrase), nil, nil)
	ck(err)
	defer encrypter.Close()
	_, err = io.Copy(encrypter, src)
	ck(err)
}

func symDecrypt(passphrase string, src io.Reader, dst io.Writer) {
	keyRing := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return []byte(passphrase), nil
	}
	md, err := openpgp.ReadMessage(src, nil, keyRing, nil)
	ck(err)
	_, err = io.Copy(dst, md.UnverifiedBody)
	ck(err)
}

func ck(err error) {
	if err != nil {
		panic("OpenPGP: " + strings.Replace(err.Error(), "openpgp: ", "", 1))
	}
}
