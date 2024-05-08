// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/system"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
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

func (pgp *suOpenPGP) Equal(other any) bool {
	return pgp == other
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
		return opgpStr(passphrase, source, symEncrypt)
	}
	return opgpFile(passphrase, source, toFile, symEncrypt)
}

var _ = staticMethod(opgp_SymmetricDecrypt,
	"(passphrase, source, toFile = false)")

func opgp_SymmetricDecrypt(passphrase, source, toFile Value) Value {
	if toFile == False {
		return opgpStr(passphrase, source, symDecrypt)
	}
	return opgpFile(passphrase, source, toFile, symDecrypt)
}

var _ = staticMethod(opgp_PublicEncrypt,
	"(publicKey, source, toFile = false)")

func opgp_PublicEncrypt(publicKey, source, toFile Value) Value {
	if toFile == False {
		return opgpStr(ToStr(publicKey), source, asymEncrypt)
	}
	return opgpFile(ToStr(publicKey), source, toFile, asymEncrypt)
}

var _ = staticMethod(opgp_PrivateDecrypt,
	"(privateKey, passphrase, source, toFile = false)")

func opgp_PrivateDecrypt(privateKey, passphrase, source, toFile Value) Value {
	kp := keyPair{privateKey: ToStr(privateKey), passphrase: ToStr(passphrase)}
	if toFile == False {
		return opgpStr(kp, source, asymDecrypt)
	}
	return opgpFile(kp, source, toFile, asymDecrypt)
}

type keyPair struct {
	privateKey string
	passphrase string
}

var _ = staticMethod(opgp_KeyGen, "(name, email, passphrase)")

func opgp_KeyGen(name, email, passphrase Value) Value {
	const rsaBits = 2048
	privateKey, err := helper.GenerateKey(ToStr(name), ToStr(email),
		[]byte(ToStr(passphrase)), "rsa", rsaBits)
	ck(err)

	keyRing, err := crypto.NewKeyFromArmoredReader(strings.NewReader(privateKey))
	ck(err)

	publicKey, err := keyRing.GetArmoredPublicKey()
	ck(err)

	ob := &SuObject{}
	ob.Set(SuStr("public"), SuStr(publicKey))
	ob.Set(SuStr("private"), SuStr(privateKey))
	return ob
}

var _ = staticMethod(opgp_KeyId, "(key)")

func opgp_KeyId(key Value) Value {
	keyOb, err := crypto.NewKeyFromArmored(ToStr(key))
	ck(err)
	return SuStr(keyOb.GetHexKeyID())
}

var _ = staticMethod(opgp_KeyEntity, "(key)")

func opgp_KeyEntity(key Value) Value {
	keyOb, err := crypto.NewKeyFromArmored(ToStr(key))
	ck(err)
	e := keyOb.GetEntity()
	for name := range e.Identities {
		return SuStr(name)
	}
	return False
}

//-------------------------------------------------------------------

type encdec[T any] func(passphrase T, src io.Reader, dst io.Writer)

func opgpStr[T any](key T, source Value, f encdec[T]) Value {
	src := strings.NewReader(ToStr(source))
	dst := new(bytes.Buffer)
	f(key, src, dst)
	return SuStr(dst.String())
}

func opgpFile[T any](key T, fromFile, toFile Value, f encdec[T]) Value {
	src, err := os.Open(ToStr(fromFile))
	ck(err)
	defer src.Close()
	to := ToStr(toFile)
	dst, err := os.CreateTemp(filepath.Dir(to), "su")
	ck(err)
	defer os.Remove(dst.Name())
	defer dst.Close()
	f(key, src, dst)
	dst.Close()
	system.RenameBak(dst.Name(), to)
	return nil
}

func symEncrypt(passphrase Value, src io.Reader, dst io.Writer) {
	encrypter, err := openpgp.SymmetricallyEncrypt(
		dst, []byte(ToStr(passphrase)), nil, nil)
	ck(err)
	defer encrypter.Close()
	_, err = io.Copy(encrypter, src)
	ck(err)
}

func symDecrypt(passphrase Value, src io.Reader, dst io.Writer) {
	keyRing := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return []byte(ToStr(passphrase)), nil
	}
	md, err := openpgp.ReadMessage(src, nil, keyRing, nil)
	ck(err)
	_, err = io.Copy(dst, md.UnverifiedBody)
	ck(err)
}

func asymEncrypt(publicKey string, src io.Reader, dst io.Writer) {
	publicKeyObj, err := crypto.NewKeyFromArmored(publicKey)
	ck(err)
	publicKeyRing, err := crypto.NewKeyRing(publicKeyObj)
	ck(err)
	encryptor, err := publicKeyRing.EncryptStreamWithCompression(dst, nil, nil)
	ck(err)
	defer encryptor.Close()
	_, err = io.Copy(encryptor, src)
	ck(err)
}

func asymDecrypt(kp keyPair, src io.Reader, dst io.Writer) {
	privateKeyObj, err := crypto.NewKeyFromArmored(kp.privateKey)
	ck(err)
	privateKeyUnlocked, err := privateKeyObj.Unlock([]byte(kp.passphrase))
	ck(err)
	defer privateKeyUnlocked.ClearPrivateParams()
	privateKeyRing, err := crypto.NewKeyRing(privateKeyUnlocked)
	ck(err)
	decryptor, err := privateKeyRing.DecryptStream(src, nil, 0)
	ck(err)
	_, err = io.Copy(dst, decryptor)
	ck(err)
}

func ck(err error) {
	if err != nil {
		panic("OpenPGP: " + strings.Replace(err.Error(), "openpgp: ", "", 1))
	}
}
