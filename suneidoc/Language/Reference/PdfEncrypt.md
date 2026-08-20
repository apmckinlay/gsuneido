<div style="float:right"><span class="builtin">Builtin</span></div>

### PdfEncrypt

Generates the pieces needed to write an AES-256 encrypted PDF.

`PdfEncrypt.KeyEntries(userPass, ownerPass, permissions = -1028) => [keyEntry: keyEntry, trailerID: trailerID, encryptionKey: encryptionKey]`
: Generates a random file encryption key and derives the encryption dictionary from it.
: **NOTE:** The value -1028 allows a user to print the document and use screen readers while strictly prohibiting them from copying text, editing content, or modifying the page structure.

`PdfEncrypt.Encrypt(data, encryptionKey) => string`
: Encrypts data with AES-256. data and encryptionKey are both hex encoded, and the result is hex encoded. Use: [string.ToHex](<String/string.ToHex.md>)

For example:

```suneido
stream = "data to be encrypted!"
enc = PdfEncrypt.KeyEntries("open me", "owner")
encrypted = PdfEncrypt.Encrypt(stream.ToHex(), enc.encryptionKey).FromHex()
```

userPass is the password required to open the document, ownerPass is the password that bypasses the permissions.