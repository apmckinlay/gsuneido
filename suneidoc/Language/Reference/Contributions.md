### Contributions

``` suneido
(name) => object
```

Builds a list of contributions from the libraries in use i.e. [Libraries](<Libraries.md>)()

Scans the libraries for records with names starting with the library name (capitalized), followed by an underscore and then the specified name.

For example:
<pre>
in onelib:
<b>Onelib_foobar</b>
123

in twolib:
<b>Twolib_foobar</b>
456

Contributions("foobar")
    => #(123, 456)
</pre>

The record name must match the name of the library that contains it. e.g. Onelib_foobar will only be found in onelib.

The library records can define any valid value e.g. numbers, strings, objects, functions, classes


See also:
[GetContributions](<GetContributions.md>),
[LastContribution](<LastContribution.md>),
[OptContribution](<OptContribution.md>),
[SoleContribution](<SoleContribution.md>),
[Plugins](<Plugins.md>)
