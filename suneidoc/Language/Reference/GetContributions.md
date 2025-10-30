### GetContributions

``` suneido
(name) => object
```

Merges the results of [Contributions](<Contributions.md>)(name)

-	The contributions must be objects or functions that return an object.
-	Unnamed members are combined into a complete list.
-	Named members whose values are objects will be merged recursively.
-	If multiple libraries contribute to the same named member, and the values are not objects, a duplicate contribution exception will be thrown. i.e. one contribution cannot "override" another one


For example:
<pre>
in onelib:
<b>Onelib_values</b>
#(1, 7) 

in twolib:
<b>Twolib_values</b>
#(9, 2)

GetContributions('values') 
    => #(1, 7, 9, 2)

in onelib:
<b>Onelib_stuff</b>
#(a: (1), b: (c: (2), d: (6)))

in twolib:
<b>Twolib_stuff</b>
#(a: (3), b: (c: (4), d: (8)))

GetContributions('stuff')
    => #(a: (1, 3), b: (c: (2, 4), d: (6, 8))
</pre>


See also:
[Contributions](<Contributions.md>),
[LastContribution](<LastContribution.md>),
[OptContribution](<OptContribution.md>),
[SoleContribution](<SoleContribution.md>),
[Plugins](<Plugins.md>)
