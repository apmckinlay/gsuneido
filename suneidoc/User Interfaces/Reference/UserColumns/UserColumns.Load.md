#### UserColumns.Load

``` suneido
(columns, title, list, reset = false, deletecol = false)
```
`columns`
: The default list of columns.

`title`
: Should be unique, the settings will be saved under this name.

`list`
: A reference to the list control.

Can be used with other list controls that support:
`list.SetColumns(columns, reset = false)`
: reset is only required if you call Load with reset

`list.SetColWidth(i, width)`
: 

Normally called from New, for example:

``` suneido
New(...)
    {
    ...
    UserColumns.Load(.columns, .Title, .list)
    }
```