### OptContribution

``` suneido
(name, def) => value
```

Gets [Contributions](<Contributions.md>)(). Throws an exception if there are more than one. Returns the single contribution, or the specified def value if there are no contributions.

If the contributions are functions, you probably want to make the default a function, for example:

``` suneido
OptContribution(name, function (@args) { <def> })(...)
```


See also:
[Contributions](<Contributions.md>),
[GetContributions](<GetContributions.md>),
[LastContribution](<LastContribution.md>),
[SoleContribution](<SoleContribution.md>),
[Plugins](<Plugins.md>)
