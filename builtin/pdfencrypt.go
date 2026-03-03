// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	. "github.com/apmckinlay/gsuneido/core"
)

type suPdfEncrypt struct {
	staticClass[suPdfEncrypt]
}

func init() {
	Global.Builtin("PdfEncrypt", &suPdfEncrypt{})
}

func (*suPdfEncrypt) String() string {
	return "PdfEncrypt /* builtin class */"
}

func (pe *suPdfEncrypt) Equal(other any) bool {
	return pe == other
}

func (*suPdfEncrypt) Lookup(_ *Thread, method string) Value {
	return pdfEncryptMethods[method]
}

var pdfEncryptMethods = methods("pe")

var _ = staticMethod(pe_Members, "()")

func pe_Members() Value {
	return pe_members
}

var pe_members = methodList(pdfEncryptMethods)

// The value -1028 allows a user to print the document and use screen readers
// while strictly prohibiting them from copying text, editing content, or modifying the page structure.
var _ = staticMethod(pe_KeyEntries, "(userPass, ownerPass, permissions = -1028)")

// pe_KeyEntries generates PDF encryption dictionary entries for AES-256 (Revision: 5, Version: 5, AESV3)
// Based on PDF 1.7 Extension Level 3 and ISO 32000-1 Algorithm 3.2a
func pe_KeyEntries(th *Thread, args []Value) Value {
	fileID := make([]byte, 16)
	if _, err := rand.Read(fileID); err != nil {
		panic("PdfEncrypt.KeyEntries: " + err.Error())
	}

	userPassBytes := []byte(ToStr(args[0]))
	ownerPassBytes := []byte(ToStr(args[1]))
	p := int32(ToInt(args[2]))

	encryptionKey := make([]byte, 32)
	if _, err := rand.Read(encryptionKey); err != nil {
		panic("PdfEncrypt.KeyEntries: " + err.Error())
	}

	uValue, ueValue := pdfComputeU(userPassBytes, encryptionKey)
	oValue, oeValue := pdfComputeO(ownerPassBytes, uValue, encryptionKey)

	permsValue := pdfComputePerms(p, encryptionKey)

	encryptDict := fmt.Sprintf(
		`<</CF<</StdCF<</AuthEvent/DocOpen/CFM/AESV3/Length 32>>>>`+
			`/Filter/Standard/Length 256/O<%x>/OE<%x>/P %d/Perms<%x>/R 5/StmF/StdCF/StrF/StdCF/U<%x>/UE<%x>/V 5>>`,
		oValue, oeValue, p, permsValue, uValue, ueValue)
	trailerID := fmt.Sprintf("/ID [<%x> <%x>]", fileID, fileID)

	ob := &SuObject{}
	ob.Set(SuStr("keyEntry"), SuStr(encryptDict))
	ob.Set(SuStr("trailerID"), SuStr(trailerID))
	ob.Set(SuStr("encryptionKey"), SuStr(hex.EncodeToString(encryptionKey)))
	return ob
}

// pdfComputeO computes the O (owner password) hash and OE (encrypted key) values
// Algorithm 3.2a steps 7-8 from ISO 32000-1
// O = SHA256(ownerPass || validationSalt || U) || validationSalt || keySalt (48 bytes)
// OE = AES-CBC encrypt(encryptionKey, key=SHA256(ownerPass || keySalt || U))
func pdfComputeO(ownerPass, uValue, encryptionKey []byte) (o, oe []byte) {
	// Generate 8-byte validation salt and 8-byte key salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic("PdfEncrypt: " + err.Error())
	}
	vs, ks := salt[:8], salt[8:]

	// O = hash || validationSalt || keySalt (48 bytes total)
	o = make([]byte, 48)
	h := sha256.New()
	h.Write(ownerPass)
	h.Write(vs)
	h.Write(uValue)
	h.Sum(o[:0])
	copy(o[32:40], vs)
	copy(o[40:48], ks)

	// OE key = SHA256(ownerPass || keySalt || U)
	h2 := sha256.New()
	h2.Write(ownerPass)
	h2.Write(ks)
	h2.Write(uValue)
	key := h2.Sum(nil)

	// OE = AES-CBC encrypt(encryptionKey)
	oe = aesEncryptCBCZeroIV(encryptionKey, key)

	return o, oe
}

// Algorithm 3.2a steps 4-6 from ISO 32000-1
// U = SHA256(userPass || validationSalt) || validationSalt || keySalt (48 bytes)
// UE = AES-CBC encrypt(encryptionKey, key=SHA256(userPass || keySalt))
func pdfComputeU(userPass, encryptionKey []byte) (u, ue []byte) {
	// Generate 8-byte validation salt and 8-byte key salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic("PdfEncrypt: " + err.Error())
	}
	vs, ks := salt[:8], salt[8:]

	// U = hash || validationSalt || keySalt (48 bytes total)
	u = make([]byte, 48)
	h := sha256.New()
	h.Write(userPass)
	h.Write(vs)
	h.Sum(u[:0]) // Write hash directly into u[0:32]
	copy(u[32:40], vs)
	copy(u[40:48], ks)

	// UE key = SHA256(userPass || keySalt)
	h2 := sha256.New()
	h2.Write(userPass)
	h2.Write(ks)
	key := h2.Sum(nil)

	// UE = AES-CBC encrypt(encryptionKey, zero IV)
	ue = aesEncryptCBCZeroIV(encryptionKey, key)

	return u, ue
}

// Algorithm 3.10 from ISO 32000-1
func pdfComputePerms(p int32, encryptionKey []byte) []byte {
	// Create 16-byte plaintext: permissions (4 bytes LE) + 0xFFFFFFFF + flags + "adb"
	b := make([]byte, 16)
	binary.LittleEndian.PutUint32(b[:4], uint32(p))

	// 'F' = encrypting metadata: no
	copy(b[4:], []byte{0xff, 0xff, 0xff, 0xff, 'F', 'a', 'd', 'b', 0xff, 0xff, 0xff, 0xff})

	// Encrypt with AES-ECB using encryption key
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		panic("PdfEncrypt: " + err.Error())
	}

	perms := make([]byte, 16)
	block.Encrypt(perms, b)

	return perms
}

// aesEncryptCBCZeroIV encrypts data using AES-256-CBC with IV
// Only for OE/UE encryption per PDF spec - each key is used once.
func aesEncryptCBCZeroIV(data, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("PdfEncrypt: " + err.Error())
	}

	iv := make([]byte, aes.BlockSize) // zero iv

	// For OE/UE, data is exactly 32 bytes (encryption key), no padding needed
	if len(data)%aes.BlockSize != 0 {
		panic("PdfEncrypt: data must be a multiple of block size")
	}

	encrypted := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, data)

	return encrypted
}

var _ = staticMethod(pe_Encrypt, "(data, encryptionKey)")

func pe_Encrypt(data, encryptionKey Value) Value {
	dataBytes, err := hex.DecodeString(ToStr(data))
	if err != nil {
		panic("PdfEncrypt.Encrypt: invalid hex data: " + err.Error())
	}

	fileKeyBytes, err := hex.DecodeString(ToStr(encryptionKey))
	if err != nil {
		panic("PdfEncrypt.Encrypt: invalid encryptionKey: " + err.Error())
	}

	if len(fileKeyBytes) != 32 {
		panic("PdfEncrypt.Encrypt: encryptionKey must be 32 bytes (256 bits) for AES-256")
	}

	str := pdfEncryptAESBytes(dataBytes, fileKeyBytes)
	return SuStr(hex.EncodeToString(str))
}

func pdfEncryptAESBytes(b, key []byte) []byte {
	n := aes.BlockSize - len(b)%aes.BlockSize
	var padding [16]byte
	for i := 0; i < n; i++ {
		padding[i] = byte(n)
	}
	b = append(b, padding[:n]...)

	data := make([]byte, aes.BlockSize+len(b))
	iv := data[:aes.BlockSize]

	_, err := rand.Read(iv)
	if err != nil {
		panic("PdfEncrypt: " + err.Error())
	}

	cb, err := aes.NewCipher(key)
	if err != nil {
		panic("PdfEncrypt: " + err.Error())
	}

	mode := cipher.NewCBCEncrypter(cb, iv)
	mode.CryptBlocks(data[aes.BlockSize:], b)

	return data
}
