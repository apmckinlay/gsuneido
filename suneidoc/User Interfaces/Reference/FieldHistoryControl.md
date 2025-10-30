### FieldHistoryControl

``` suneido
(width = 20, mandatory = false, selectFirst = false,
    font = "", trim = true)
```

Derived from [ChooseListControl](<ChooseListControl.md>). Used like a [FieldControl](<FieldControl.md>), except that entries are added to the list, so that they can be chosen later rather than re-entered.

![](<../../res/ChooseListControl1.gif>)

The lists are stored in Suneido.FieldHistory under the field name, e.g. Suneido.FieldHistory.FindInName.

**Note**: Currently, the history is not saved when you exit from Suneido.

If **trim** is true, leading and trailing whitespace will be removed.