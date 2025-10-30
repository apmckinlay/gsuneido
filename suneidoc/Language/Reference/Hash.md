### Hash

``` suneido
(value) => number
```

Returns a hash code for the value. This is the same hash code that Suneido uses internally for hash tables in e.g. Objects.

Hash codes are **not** unique. If two values have different hash codes then they are guaranteed to be different values, but if they have the same hash code they are not guaranteed to be the same value.

Hash is for specialized uses, normal code should not need it.

See also: [Same?](<Same?.md>)