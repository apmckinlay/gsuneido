### RecordFormat

``` suneido
( record, fields = False )
```

Derived from
[GridFormat](<GridFormat>).

Takes a record and a list of fields to print. 
Prints the Prompt for each field beside its value.
If the fields argument is not supplied, then all the fields are printed
(in no particular order).

For example:

``` suneido
record = #(name: 'Fred', cell: '222-1234', email: 'fred@mail.com')
Params.On_Preview(Object('Record', record, fields: #(name, cell, email)))
```

Would produce:

![](<../../res/PrintRecord.gif>)