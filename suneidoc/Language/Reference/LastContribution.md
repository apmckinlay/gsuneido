### LastContribution

``` suneido
(name) => value
```

Gets [Contributions](<Contributions.md>)(). Throws an exception if there are no contributions. Returns the last contribution (in [Libraries](<Libraries.md>)() order).

For example, if the libraries in use are stdlib, onelib, and twolib (in that order) and onelib contains Onelib_foobar and twolib contains Twolib_foobar, then LastContribution("foobar") will return the value of Two_foobar.

Useful when code requires an externally defined service. It can provide a default implementation in it's own library, which can then be overridden if desired.


See also:
[Contributions](<Contributions.md>),
[GetContributions](<GetContributions.md>),
[OptContribution](<OptContribution.md>),
[SoleContribution](<SoleContribution.md>),
[Plugins](<Plugins.md>)
