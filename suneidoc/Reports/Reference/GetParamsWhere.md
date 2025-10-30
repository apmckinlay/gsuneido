#### GetParamsWhere

``` suneido
(field, ..., data: object) => string
```

Returns a string containing the "where" portion of a query.

GetParamsWhere is normally called via [report.GetParamsWhere](<Report/report.GetParamsWhere.md>) in the Query method of a [QueryFormat](<QueryFormat.md>).  The data normally comes from [ParamsSelectControls](<../../User Interfaces/Reference/ParamsSelectControl.md>), whose values consist of an object with `operation`, `value`, and `value2` members.

For example:

``` suneido
GetParamsWhere("first_name", "last_name", data: #{
    first_name: #(operation: "equals", value: "Robert", value2: ""),
    last_name: #(operation: "equals", value: "Johnson", value2: "")})
        => 'where first_name="Robert" and last_name="Johnson"'
```

See also:
[Params](<Params.md>),
[Going Further - Report Parameters](<../../Going Further/Report Parameters.md>)