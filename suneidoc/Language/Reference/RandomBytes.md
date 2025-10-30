<div style="float:right"><span class="builtin">Builtin</span></div>

### RandomBytes

``` suneido
(nbytes) => string
```

Returns a string of cryptographically random bytes of length **nbytes**.

**nbytes** must be in the range 0 to 128.
The size is limited because there is a limited amount of entropy available.
This function should not be used to generate large amounts of random data.

See also: [Random](<Random.md>)