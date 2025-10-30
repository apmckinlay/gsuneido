<div style="float:right"><span class="builtin">Builtin</span></div>

### Sha1

``` suneido
(string) => number
() => instance
```

If called with a string it simply returns the 20 byte SHA1 checksum of the string.

For example:

``` suneido
Sha1("hello world").ToHex()
    => "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"
```

To accumulate a checksum incrementally, you can create an instance and update it. For example:

``` suneido
cksum = Sha1()
cksum.Update("hello ")
cksum.Update("world")
cksum.Value().ToHex()
    => "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"
```

**Note**: Once Value() has been called, you can no longer continue to Update the checksum.


See also:
[Adler32](<Adler32.md>),
[Md5](<Md5.md>), <a href="/suneidoc/Language/Reference/ResourceCounts">ResourceCounts</a>
