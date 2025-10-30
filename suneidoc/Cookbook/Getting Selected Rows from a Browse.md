## Getting Selected Rows from a Browse

**Category:** User Interface

**Problem:**

How to retrieve the selected row(s) from a BrowseControl

**Ingredients:**

BrowseControl

**Recipe:**

To get the selected row(s)

``` suneido
sel = browse.GetSelection() // this returns an object containing the indexes of the selected rows
```

To get the data in a selected row

``` suneido
row = sel[0] // this will be the index of the first row selected if there are multiple rows selected
data = browse.GetRow(row) // data will contain all the values from all the fields on that row
```
**Example:**
``` suneido
Controller
    {
    Controls: #(Vert
        (Browse 'tables')
        (Button 'Get Record'))
    On_Get_Record()
        {
        browse = .Vert.Browse
        selected = browse.GetSelection()
        if selected.Empty?()
            {
            Alert('You have not selected a row')
            return
            }
        if selected.Size() > 1
            {
            Alert('You can only select one row')
            return
            }
        index = selected[0]
        rec = browse.GetRow(index)
        Alert('You have selected table ' $ Display(rec.tablename))
        }
    }
```