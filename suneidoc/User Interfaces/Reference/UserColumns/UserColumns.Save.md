#### UserColumns.Save

``` suneido
(title, list)
```
`title`
: Should be unique, the settings will be saved under this name.

`list`
: A reference to the list control.

Can be used with other list controls that support:
`list.HeaderChanged?()`
: optional - if implemented, the columns will only be saved if HeaderChanged? returns true

`list.GetColumns()`
: 

`list.GetColWidth(i)`
: 

Normally called from Destroy, for example:

``` suneido
Destroy()
    {
    UserColumns.Save(.Title, .list)
    super.Destroy()
    }
```