### SoleContribution

``` suneido
(name) => value
```

Gets [Contributions](<Contributions.md>)(). Throws an exception if there are no contributions or if there are more than one. Returns the single contribution.

Useful when code requires an externally defined service but you don't want to hard code the dependency. You could just reference an undefined name, but this does not give a helpful error message and doesn't document what is going on.


See also:
[Contributions](<Contributions.md>),
[GetContributions](<GetContributions.md>),
[LastContribution](<LastContribution.md>),
[OptContribution](<OptContribution.md>),
[Plugins](<Plugins.md>)
