<div style="float:right"><span class="builtin">Builtin</span></div>

### Bloom

``` suneido
(n, p) => bloom
```

Creates a [Bloom filter](<https://en.wikipedia.org/wiki/Bloom_filter>) for an estimated n values with a probability p of false positives. For example, if p is .01 there is a 1/100 chance of false positives. The size of the Bloom filter will increase with larger n or smaller p. The entire size if allocated when the Bloom filter is created.

Methods:
`Add(value)`
: Adds a value.

`Test(value) => true or false`
: Returns false if the value has definitely not been added and true if it has **probably** been added.

`Size() => bytes`
: Returns the number of bytes used by the Bloom filter.