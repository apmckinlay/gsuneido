<div style="float:right"><span class="builtin">Builtin</span></div>

### Argon2id

``` suneido
(password, salt) => string
```

Derives a 32-byte key from the password and salt using Argon2id,
the hybrid version of Argon2 that combines data-independent and data-dependent memory access.

Uses time cost 1, 64 MiB memory, and 1 thread — the second recommended option from RFC 9106 Section 4.

For example:

``` suneido
Argon2id("password", "somesalt").ToHex()
    => "d15801...6c2a"
```

Argon2id is intended for password hashing. The salt should be unique per password, typically 16 or more random bytes.


See also:
[Adler32](<Adler32.md>),
[Md5](<Md5.md>),
[Sha1](<Sha1.md>),
[Sha256](<Sha256.md>)
