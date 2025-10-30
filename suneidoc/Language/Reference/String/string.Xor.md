#### string.Xor

``` suneido
(string) => string
```

str.Xor(key) exclusive or's each byte of str with the corresponding byte of key.

If key is shorter than str, it will be effectively repeated.

This can be useful for obfuscating strings but should not be considered "secure".