### MemoizeSingle

Abstract base class for caching the results of a function with no arguments.

Derived classes must define Func.

For example, LibraryTables is defined with MemoizeSingle since it is potentially slow to derive.

``` suneido
MemoizeSingle
    {
    Func()
        {
        query = 'columns
                    summarize table, list column
                    where HasLibraryColumns?(list_column)
                join tables'
        return QueryList(query, 'tablename')
        }
    }
```

The cached result is stored in the global Suneido object.

MemoizeSingle also defines a ResetCache method that will reset the associated cached value, and a ClearAll method that will the will reset all MemoizeSingle (ResetCaches uses this)


See also:
[Compose](<Compose.md>),
[Curry](<Curry.md>),
[Memoize](<Memoize.md>),
[object.Any?](<Object/object.Any?.md>),
[object.Drop](<Object/object.Drop.md>),
[object.Every?](<Object/object.Every?.md>),
[object.Filter](<Object/object.Filter.md>),
[object.FlatMap](<Object/object.FlatMap.md>),
[object.Fold](<Object/object.Fold.md>),
[object.Map](<Object/object.Map.md>),
[object.Map!](<Object/object.Map!.md>),
[object.Map2](<Object/object.Map2.md>),
[object.Nth](<Object/object.Nth.md>),
object.PrevSeq,
[object.Reduce](<Object/object.Reduce.md>),
[object.Take](<Object/object.Take.md>),
[object.Zip](<Object/object.Zip.md>)
