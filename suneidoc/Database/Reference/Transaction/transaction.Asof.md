<div style="float:right"><span class="builtin">Builtin</span></div>

#### transaction.Asof

``` suneido
(asof = false) => date or false
```

Asof allows setting and getting the date-time that the transaction is "as of".

**Note**: The database only saves its state periodically e.g. every minute. And it only saves the state if there have been changes. Multiple transaction updates may be included in a saved state. Only the saved states can be viewed with Asof.
`Asof() => date-time or false`
: Gets the current as-of date, or false if there isn't one.

`Asof(datetime) => date`
: Sets the transaction's as-of to the closest state less than or equal to the specified date-time. It returns the date-time of the state.

`Asof(+1) or Asof(-1) => date-time or false`
: Sets the transaction's as-of to the next (+1) or previous (-1) state and returns the date-time of the state, or false if it hits the beginning or end.

For example, to find previous versions of a library record:

``` suneido
text = Query1('stdlib', :name).text
Transaction(read:)
    {|t|
    while false isnt asof = t.Asof(-1)
        {
        x = t.Query1('stdlib', :name)
        if x.text isnt text
            {
            Print(x.lib_modified)
            text = x.text
            }
        }
    }
```

**Note**: This is intended for debugging purposes. It should not be used in actual application code. It depends on internal database features that could change at any time.

**Note**: Asof is relatively slow. It potentially scans the entire database file, which could be multiple gigabytes.

**Note**: The information used for Asof is removed when the database is compacted.

See also: history virtual table