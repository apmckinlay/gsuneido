### InfoFormat

``` suneido
( prefix = '', data = False )
```

Prints any "info" fields (with the specified prefix) in the data
vertically, in a column on a report.

For example:

``` suneido
record = #(info2: 'Work: 222-7856', info3: 'Home: 123-4567', info5: 'Email: joe@mail.com')
Params.On_Preview(Object('Info' data: record))
```

Would produce:

![](<../../res/InfoFormat.gif>)

See also:
[InfoControl](<../../User Interfaces/Reference/InfoControl.md>)