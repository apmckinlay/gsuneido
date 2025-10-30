#### string.SplitFixedLength

``` suneido
(map) => record
```

SplitFixedLength extracts the values from a fixed length format string using the specified **map**.

**map** should be an object with field names for members. The values should be objects containing the start position "pos" and length "len".

For example:

``` suneido
"abcdefghijklmnopqrstuvwxyz".SplitFixedLength(
    #(field3: (pos:2, len: 5) field4: (pos: 7, len: 2) field5: (pos: 9, len: 12)))
        => [field3: 'cdefg' field4: 'hi' field5: 'jklmnopqrstu']
```