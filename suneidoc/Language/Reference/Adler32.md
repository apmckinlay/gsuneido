<div style="float:right"><span class="builtin">Builtin</span></div>

### Adler32

``` suneido
(string) => number
() => instance
```

If called with a string it simply returns the Adler32 checksum of the string.

For example:

``` suneido
Adler32("hello world")
    => 436929629
```

To accumulate a checksum incrementally, you can create an instance and update it. For example:

``` suneido
cksum = Adler32()
cksum.Update("hello ")
cksum.Update("world")
cksum.Value()
    => 436929629
```


See also:
[Md5](<Md5.md>),
[Sha1](<Sha1.md>)
