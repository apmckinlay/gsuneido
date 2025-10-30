<div style="float:right"><span class="builtin">Builtin</span></div>

### OpenPGP
`OpenPGP.SymmetricEncrypt(passphrase, string) => string`
: Encrypts string and returns the result as a (binary) string.

`OpenPGP.SymmetricEncrypt(passphrase, fromFile, toFile)`
: Encrypts the contents of fromFile and stores it in toFile.

`OpenPGP.SymmetricDecrypt(passphrase, string) => string`
: Decrypts string and returns the result as a string.

`OpenPGP.SymmetricDecrypt(passphrase, fromFile, toFile)`
: Decrypts the contents of fromFile and stores it in toFile.

`OpenPGP.KeyGen(name, email, passphrase) => Object(public: string, private: string)`
: Generates a public/private key pair with RSA bits of 2048.

`OpenPGP.PublicEncrypt(publicKey, string) => string`
: Encrypts string and returns the result as a (binary) string.

`OpenPGP.PublicEncrypt(publicKey, fromFile, toFile)`
: Encrypts the contents of fromFile and stores it in toFile.

`OpenPGP.PrivateDecrypt(privateKey, passphrase, string) => string`
: Decrypts string and returns the result as a string.

`OpenPGP.PrivateDecrypt(privateKey, passphrase, fromFile, toFile)`
: Decrypts the contents of fromFile and stores it in toFile.

`OpenPGP.KeyId(key)`
: Returns the key id as a hex encoded string.

`OpenPGP.KeyEntity(key)`
: Returns the name of the first key entity.

**Note:** With potentially large files, it is better to use the file versions, rather than dealing with large amounts of data in memory.

These methods should be interoperable with other OpenPGP compatible software such as GnuPG (gpg).

Symmetric encrypt uses the defaults - as of Sept. 2022 this was AES-256 and SHA1.

Asymmetric encrypt compresses before encryption (as of Feb. 2024 with zlib)