### Subscript

Subscript operations [...] can be applied to strings and objects.

`ob[m]` returns the value of the specified member of the object. If m is not a member of the object it will throw an exception of "member not found..."

``` suneido
ob = #(12, 34, 56)
ob[0]
    => 12
ob[2]
    => 56
ob[9]
    => member not found 9
```

`str[i]` returns a one character string containing the character at the i'th position. Negative indexes are taken relative to the end of the string.  Indexes outside range of the string will return an empty string ("").

``` suneido
s = "hello"
s[0]
    => "h"
s[4]
    => "o"
s[5]
    => ""
s[-1]
    => "o"
s[-5]
    => "h"
```

Range subscripts may also be used.

`[from .. to]` returns the subsequence starting at `from` and ending at, but not including `to`.

`s[from :: length]` returns the subsequence starting at `from` with the given `length`.

If `from` or `to` are negative they are taken as relative to the end, with a minimum of 0. (Note: a negative length is **not** taken as end relative, use '..' for this situation)

The subsequence is limited to the range of the string. Indexes out of range will not throw an exception.

If from is omitted, the subsequence will start at the beginning. If `to` or `length` are ommitted, the subsequence will extend to the end.

For example:

``` suneido
s = "foobar"
s[0 .. 3]
    => "foo"
s[.. 3]
    => "foo"
s[2 .. 4]
    => "ob"
s[1 .. -1]
    => "ooba"
s[3 .. 6]
    => "bar"
s[3 ..]
    => "bar"
s[2 :: 2]
    => "ob"
s[3 :: 3]
    => "bar"
s[:: 2]
    => "fo"
s[1 :: -1]
    => ""

ob = #(a, b, c, d)
ob[1 .. -1]
    => #(b, c)
ob[:: 2]
    => #(a, b)
ob[2 ..]
    => #(c, d)
ob[1 :: -1]
    => #()
```

[.. n] is the same as [:: n] (both give the first n characters or values).

**Note:** When using literal numbers, spaces are required around the ".." to prevent the periods being treated as decimal points.

**Note:** For objects, ranges only apply to the initial un-named members part of the object. For example:

``` suneido
#(a, b, 7: c, 9: d)[0 .. 9]
    => #(a, b) // members 7 and 9 not included
```