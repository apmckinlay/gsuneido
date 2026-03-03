// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestPdfEncrypt_KeyEntries(t *testing.T) {
	th := &Thread{}

	// Test with default permissions
	result := pe_KeyEntries(th, []Value{SuStr("userpass"), SuStr("ownerpass"), IntVal(-1028)})
	ob := result.(*SuObject)

	// Verify the result has the expected keys
	keyEntry := ob.Get(th, SuStr("keyEntry"))
	trailerID := ob.Get(th, SuStr("trailerID"))
	encryptionKey := ob.Get(th, SuStr("encryptionKey"))

	assert.T(t).That(keyEntry != nil)
	assert.T(t).That(trailerID != nil)
	assert.T(t).That(encryptionKey != nil)

	// Verify keyEntry contains expected PDF encryption dictionary fields
	keyEntryStr := ToStr(keyEntry)
	assert.T(t).That(bytes.Contains([]byte(keyEntryStr), []byte("/Filter/Standard")))
	assert.T(t).That(bytes.Contains([]byte(keyEntryStr), []byte("/V 5")))
	assert.T(t).That(bytes.Contains([]byte(keyEntryStr), []byte("/R 5")))
	assert.T(t).That(bytes.Contains([]byte(keyEntryStr), []byte("/Length 256")))

	// Verify encryption key is 64 hex characters (32 bytes)
	encKeyStr := ToStr(encryptionKey)
	assert.T(t).This(len(encKeyStr)).Is(64)

	// Verify trailerID contains /ID
	trailerIDStr := ToStr(trailerID)
	assert.T(t).That(bytes.Contains([]byte(trailerIDStr), []byte("/ID")))
}

func TestPdfEncrypt_KeyEntries_CustomPermissions(t *testing.T) {
	th := &Thread{}

	// Test with custom permissions
	result := pe_KeyEntries(th, []Value{SuStr("userpass"), SuStr("ownerpass"), IntVal(0)})
	ob := result.(*SuObject)

	keyEntry := ob.Get(th, SuStr("keyEntry"))
	keyEntryStr := ToStr(keyEntry)

	// Verify P value is set
	assert.T(t).That(bytes.Contains([]byte(keyEntryStr), []byte("/P 0")))
}

func TestPdfComputeU(t *testing.T) {
	userPass := []byte("testpassword")
	encryptionKey := make([]byte, 32)
	for i := range encryptionKey {
		encryptionKey[i] = byte(i)
	}

	u, ue := pdfComputeU(userPass, encryptionKey)

	// U should be 48 bytes: hash (32) + validationSalt (8) + keySalt (8)
	assert.T(t).This(len(u)).Is(48)
	// UE should be 32 bytes (encrypted encryption key)
	assert.T(t).This(len(ue)).Is(32)

	// Verify U hash is correct
	vs := u[32:40]
	ks := u[40:48]
	expectedHash := sha256.Sum256(append(userPass, vs...))
	assert.T(t).This(u[:32]).Is(expectedHash[:])

	// Verify UE can be decrypted with correct key
	ueKey := sha256.Sum256(append(userPass, ks...))
	block, _ := aes.NewCipher(ueKey[:])
	decrypted := make([]byte, 32)
	mode := cipher.NewCBCDecrypter(block, make([]byte, aes.BlockSize))
	mode.CryptBlocks(decrypted, ue)
	assert.T(t).This(decrypted).Is(encryptionKey)
}

func TestPdfComputeO(t *testing.T) {
	ownerPass := []byte("ownerpassword")
	userPass := []byte("userpass")
	encryptionKey := make([]byte, 32)
	for i := range encryptionKey {
		encryptionKey[i] = byte(i)
	}

	// First compute U since O depends on it
	u, ue := pdfComputeU(userPass, encryptionKey)

	o, oe := pdfComputeO(ownerPass, u, encryptionKey)

	// O should be 48 bytes: hash (32) + validationSalt (8) + keySalt (8)
	assert.T(t).This(len(o)).Is(48)
	// OE should be 32 bytes (encrypted encryption key)
	assert.T(t).This(len(oe)).Is(32)

	// Extract salts from computed O
	vs := o[32:40]
	ks := o[40:48]

	// Verify O hash is correct: SHA256(ownerPass || validationSalt || U)
	// U is the full 48-byte value
	expectedHash := sha256.Sum256(append(append(ownerPass, vs...), u...))
	assert.T(t).This(o[:32]).Is(expectedHash[:])

	// Verify OE can be decrypted with correct key
	// OE key = SHA256(ownerPass || keySalt || U)
	// U is the full 48-byte value
	oeKey := sha256.Sum256(append(append(ownerPass, ks...), u...))
	block, _ := aes.NewCipher(oeKey[:])
	decrypted := make([]byte, 32)
	iv := make([]byte, aes.BlockSize)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, oe)
	assert.T(t).This(decrypted).Is(encryptionKey)

	// Additional verification: decrypt UE to recover the encryption key
	// and verify we can use it to validate O
	ueKey := sha256.Sum256(append(userPass, u[40:48]...))
	block2, _ := aes.NewCipher(ueKey[:])
	recoveredKey := make([]byte, 32)
	mode2 := cipher.NewCBCDecrypter(block2, iv)
	mode2.CryptBlocks(recoveredKey, ue)
	assert.T(t).This(recoveredKey).Is(encryptionKey)
}

func TestPdfComputePerms(t *testing.T) {
	encryptionKey := make([]byte, 32)
	for i := range encryptionKey {
		encryptionKey[i] = byte(i)
	}

	perms := pdfComputePerms(-1028, encryptionKey)

	// Perms should be 16 bytes
	assert.T(t).This(len(perms)).Is(16)

	// Decrypt and verify permissions
	block, _ := aes.NewCipher(encryptionKey)
	decrypted := make([]byte, 16)
	block.Decrypt(decrypted, perms)

	// First 4 bytes should be permissions in little-endian
	p := uint32(decrypted[0]) | uint32(decrypted[1])<<8 | uint32(decrypted[2])<<16 | uint32(decrypted[3])<<24
	assert.T(t).This(int32(p)).Is(int32(-1028))

	// Bytes 4-7 should be 0xffffffff
	assert.T(t).This(decrypted[4:8]).Is([]byte{0xff, 0xff, 0xff, 0xff})

	// Byte 8 should be 'F' (encrypt metadata: no)
	assert.T(t).This(decrypted[8]).Is(byte('F'))

	// Bytes 9-11 should be "adb"
	assert.T(t).This(decrypted[9:12]).Is([]byte("adb"))
}

func TestPdfComputePerms_ZeroPermissions(t *testing.T) {
	encryptionKey := make([]byte, 32)
	for i := range encryptionKey {
		encryptionKey[i] = byte(i)
	}

	perms := pdfComputePerms(0, encryptionKey)

	block, _ := aes.NewCipher(encryptionKey)
	decrypted := make([]byte, 16)
	block.Decrypt(decrypted, perms)

	p := uint32(decrypted[0]) | uint32(decrypted[1])<<8 | uint32(decrypted[2])<<16 | uint32(decrypted[3])<<24
	assert.T(t).This(p).Is(uint32(0))
}

func TestAesEncryptCBC(t *testing.T) {
	// Test with 32 bytes of data (one block past IV)
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	data := make([]byte, 32)
	for i := range data {
		data[i] = byte(i * 2)
	}

	encrypted := aesEncryptCBCZeroIV(data, key)

	// Encrypted data should be same length as input
	assert.T(t).This(len(encrypted)).Is(32)

	// Verify we can decrypt it
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, 32)
	mode := cipher.NewCBCDecrypter(block, make([]byte, aes.BlockSize))
	mode.CryptBlocks(decrypted, encrypted)
	assert.T(t).This(decrypted).Is(data)
}

func TestAesEncryptCBC_MultipleBlocks(t *testing.T) {
	key := make([]byte, 32)
	data := make([]byte, 64) // Two blocks

	for i := range data {
		data[i] = byte(i)
	}

	encrypted := aesEncryptCBCZeroIV(data, key)
	assert.T(t).This(len(encrypted)).Is(64)

	// Verify decryption
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, 64)
	mode := cipher.NewCBCDecrypter(block, make([]byte, aes.BlockSize))
	mode.CryptBlocks(decrypted, encrypted)
	assert.T(t).This(decrypted).Is(data)
}

func TestPdfEncrypt_Encrypt(t *testing.T) {
	// Test data and key
	data := []byte("Hello, World!")
	dataHex := hex.EncodeToString(data)
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	keyHex := hex.EncodeToString(key)

	result := pe_Encrypt(SuStr(dataHex), SuStr(keyHex))
	resultStr := ToStr(result)

	// Result should be hex encoded
	encrypted, err := hex.DecodeString(resultStr)
	assert.T(t).This(err).Is(nil)

	// Encrypted data should be: IV (16 bytes) + padded data
	// Original data is 13 bytes, padded to 16 bytes
	assert.T(t).This(len(encrypted)).Is(32) // 16 (IV) + 16 (padded data)

	// Verify we can decrypt it
	iv := encrypted[:16]
	ciphertext := encrypted[16:]

	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, ciphertext)

	// Remove PKCS7 padding
	paddingLen := int(decrypted[len(decrypted)-1])
	decrypted = decrypted[:len(decrypted)-paddingLen]
	assert.T(t).This(decrypted).Is(data)
}

func TestPdfEncrypt_Encrypt_ExactBlockSize(t *testing.T) {
	// Test with data that is exactly one block size (16 bytes)
	data := make([]byte, 16)
	for i := range data {
		data[i] = byte(i)
	}
	dataHex := hex.EncodeToString(data)
	key := make([]byte, 32)
	keyHex := hex.EncodeToString(key)

	result := pe_Encrypt(SuStr(dataHex), SuStr(keyHex))
	resultStr := ToStr(result)

	encrypted, _ := hex.DecodeString(resultStr)

	// Should add a full block of padding (16 bytes)
	assert.T(t).This(len(encrypted)).Is(48) // 16 (IV) + 32 (padded data)
}

func TestPdfEncrypt_Encrypt_InvalidHex(t *testing.T) {
	assert.T(t).This(func() {
		pe_Encrypt(SuStr("not valid hex"), SuStr("0000000000000000000000000000000000000000000000000000000000000000"))
	}).Panics("invalid hex")
}

func TestPdfEncrypt_Encrypt_InvalidKeyLength(t *testing.T) {
	dataHex := hex.EncodeToString([]byte("test"))
	shortKeyHex := hex.EncodeToString([]byte("short"))

	assert.T(t).This(func() {
		pe_Encrypt(SuStr(dataHex), SuStr(shortKeyHex))
	}).Panics("must be 32 bytes")
}

func TestPdfEncryptAESBytes(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	testCases := []struct {
		dataLen     int
		expectedLen int // IV + padded data
	}{
		{1, 32},    // 1 byte + 15 padding = 16, plus 16 byte IV
		{15, 32},   // 15 bytes + 1 padding = 16, plus 16 byte IV
		{16, 48},   // 16 bytes + 16 padding = 32, plus 16 byte IV
		{31, 48},   // 31 bytes + 1 padding = 32, plus 16 byte IV
		{32, 64},   // 32 bytes + 16 padding = 48, plus 16 byte IV
		{100, 128}, // 100 bytes + 12 padding = 112, plus 16 byte IV
	}

	for _, tc := range testCases {
		data := make([]byte, tc.dataLen)
		for i := range data {
			data[i] = byte(i)
		}

		encrypted := pdfEncryptAESBytes(data, key)
		assert.T(t).This(len(encrypted)).Is(tc.expectedLen)

		// Verify decryption
		iv := encrypted[:16]
		ciphertext := encrypted[16:]

		block, _ := aes.NewCipher(key)
		decrypted := make([]byte, len(ciphertext))
		mode := cipher.NewCBCDecrypter(block, iv)
		mode.CryptBlocks(decrypted, ciphertext)

		// Remove PKCS7 padding
		paddingLen := int(decrypted[len(decrypted)-1])
		decrypted = decrypted[:len(decrypted)-paddingLen]
		assert.T(t).This(decrypted).Is(data)
	}
}
