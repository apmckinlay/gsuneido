### Opt

``` suneido
(@strings)
```

If any of the strings are "", the result is "". Otherwise the result is simply the concatenation of the strings (with no separator).

This is useful for adding "optional" (thus the name "Opt") strings with separators or delimiters. For example:

``` suneido
Opt("(", area, ") ") $ phone
```

If area is "" this will just be "phone", otherwise it will be "(area) phone"

Or, to add an optional salutation:

``` suneido
Opt(salutation, " ") $ name
```

See also:
[object.Join](<Object/object.Join.md>)