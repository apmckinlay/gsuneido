### BuildQueryWhere

``` suneido
( restrictions ) => string
```

Returns a "where" clause for a query from the given restrictions.  Restrictions is a
nested object, the inner objects each containing three unnamed members: one for the
field, one for the operator, and one for the value (in that order).

For example:

``` suneido
BuildQueryWhere(#((age, ">", 24), (city, "==", "Saskatoon")))
    => 'where age>24 and city=="Saskatoon"'
```