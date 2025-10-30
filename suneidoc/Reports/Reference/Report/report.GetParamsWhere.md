#### report.GetParamsWhere

``` suneido
(field, ...) => string
```

Calls [GetParamsWhere](<../GetParamsWhere.md>), passing .Params as the data.

For example:

``` suneido
...
Params:
    #(Form
        (ParamsSelect firstname)
        (ParamsSelect lastname))
QueryFormat
    {
    Query()
        {
        return 'mycontacts' $ _report.GetParamsWhere('firstname, lastname')
        }
...
```

See also:
[Params](<../Params.md>),
[Going Further - Report Parameters](<../../../Going Further/Report Parameters.md>)