<div style="float:right"><span class="builtin">Builtin</span></div>

#### Database.CurrentSize

``` suneido
() => number
```

Returns the current size of the database in bytes.

Note: An open database file will be larger than this size, due to the margin for growth.

For example:

``` suneido
Database.CurrentSize()
    => 10601476
```