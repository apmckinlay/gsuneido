<div style="float:right"><span class="builtin">Builtin</span></div>

### Md5

``` suneido
(string) => string
() => instance
```

If called with a string it simply returns the 16 byte MD5 hash of the string.

For example:

``` suneido
Md5("hello world").ToHex()
    => "5eb63bbbe01eeed093cb22bb8f5acdc3"
```

To accumulate a checksum incrementally, you can create an instance and update it. For example:

``` suneido
cksum = Md5()
cksum.Update("hello ")
cksum.Update("world")
cksum.Value().ToHex()
    => "5eb63bbbe01eeed093cb22bb8f5acdc3"
```

**Note**: Once Value() has been called, you can no longer continue to Update the checksum.


See also:
[Adler32](<Adler32.md>),
[Sha1](<Sha1.md>), <a href="/suneidoc/Language/Reference/ResourceCounts">ResourceCounts</a>
