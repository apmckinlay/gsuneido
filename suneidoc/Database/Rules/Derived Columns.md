### Derived Columns

When a column named is created with a capitalized name, it is a derived (calculated) column.  
It is not physically stored in the database.  
These columns can be used for searching and sorting, but they cannot be used as index or key columns.

A function called "Rule_" $ fieldname must exist to calculate the column value.  
It is called as if it were a method of the record so it may access the other columns of the record as dot variables.

For example:

``` suneido
create mytable (item quantity rate Total) key (item)

Rule_total
function()
    { return .quantity * .rate; }
```

The total column of mytable will now be automatically calculated by multiplying the quantity column times the rate column.  The total column can be used like other columns:

``` suneido
mytable where total > 100 sort total
```