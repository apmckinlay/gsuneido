#### object.JoinCSV

``` suneido
(fields = false) => string
```

JoinCSV combines the values from an object into a comma separated string.

-	strings are contained in double quotes
-	double quotes within strings are replaced with two double quotes


If **fields** is specified then only these fields will be placed into the string,
in the order specified.
If the object contains named members,
you will need to specify fields or else the order of the members will be unpredictable.

For example:

``` suneido
#(123, '35" of string').JoinCSV()
    => '123,"35"" of string"'
```

Or with named members:

``` suneido
#(name: "Andrew McKinlay", age: 41).JoinCSV(#("name", "age"))
    => '"Andrew McKinlay",41'
```

See also:
[string.SplitCSV](<../String/string.SplitCSV.md>)