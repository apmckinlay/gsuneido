#### number.ToWordsSimple

``` suneido
() => string
```

Returns the number converted to a string of words.

For example:

``` suneido
0.ToWordsSimple() => "* zero *"
8.ToWordsSimple() => "* eight *"
18.ToWordsSimple() => "* one eight *"
(-12.34).ToWordsSimple() => "* minus one two point three four *"
```

**Note:** Any fractional portion (e.g. cents) is ignored.

See also:
[number.ToWords](<number.ToWords.md>)