### FormatQuery

``` suneido
(query) => string
```

Formats the query in a standard one with one operation per line and nesting shown by indentation. It also shows query warnings as comments (project not unique, union not disjoint, join many to many)