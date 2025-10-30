## Using Excel from Suneido

by Claudio Mascioni

**Category:** Coding

**Problem**

Read and write cell values from/to an Excel spreadsheet from Suneido.

**Ingredients**

COMobject

**Recipe**

To open an Excel sheet in an Excel program window from Suneido:

``` suneido
excelapp = COMobject("Excel.Application")

excelapp.Visible = true
    // true  = display the Excel window
    // false = hide the Excel window

wrkbks = excelapp.Workbooks.Open("c:\\test.xls")  // xls file to open

wrkbks.Release()
excelapp.Release()   // to free the excel object memory
```

To import a cell value from an Excel sheet:

``` suneido
excelapp = COMobject("Excel.Application")
excelapp.Visible = false    // hide the Excel window
wrkbks = excelapp.Workbooks.Open('c:/test.xls')
cellvalue = excelapp.Cells(2,2).Value // get a value from a xls cell
wrkbks.Close    // close the excel Workbooks
Print(Display(cellvalue))
wrkbks.Release()
excelapp.Release()
```

To write a value into a cell of an existing Excell sheet and save the file:

``` suneido
excelapp = COMobject("Excel.Application")
excelapp.Visible = true;    // display the sheet with Excel program
wrkbks = excelapp.Workbooks.Open("c:\\test.xls")
excelapp.Cells(2,2).Value = Display(Timestamp()) // write something in the 2,2 cell
wrkbks.Save
wrkbks.Close
wrkbks.Release()
excelapp.Release()
```

To write a value into a cell of a new Excel sheet and display the sheet:

``` suneido
excelapp = COMobject("Excel.Application")
excelapp.Visible = true;     // display the sheet with Excel program
wrkbks = excelapp.Workbooks.Add()
excelapp.Cells(2,2).Value = Display(Timestamp()) // write something
wrkbks.Release()
excelapp.Release()
```

To write a value in a cell of a new Excell sheet and save in a new file:

``` suneido
excelapp = COMobject("Excel.Application")
excelapp.Visible = false;   // hide the Excel window
wrkbks = excelapp.Workbooks.Add()
excelapp.Cells(2,2).Value = Display(Timestamp()) // write something
wrkbks.SaveAs("c:\\testnew.xls")
wrkbks.Close
wrkbks.Release()
excelapp.Release()
```

To change the name of the sheets displayed in the tab at the end of the sheet and put a value in each sheet:

``` suneido
excelapp = COMobject("Excel.Application")
excelapp.Visible = true;
wrkbks = excelapp.Workbooks.Add()
//
sheet1 = wrkbks.WorkSheets(1)
sheet1.Name = 'nameofsheet1'             // change the tab name
sheet1.Cells(1,1).Value = 'valuesheet1'  // write something
//
sheet2 = wrkbks.WorkSheets(2)
sheet2.Name = 'nameofsheet2'
sheet2.Cells(2,2).Value = 'valuesheet2'
//
sheet3 = wrkbks.WorkSheets(3)
sheet3.Name = 'nameofsheet3'
sheet3.Cells(3,3).Value = 'valuesheet3'
//
wrkbks.Release()
sheet1.Release()
sheet2.Release()
sheet3.Release()
excelapp.Release()
```

**Bugs**

There is a problem with these examples - Excel is not released from memory. If anyone has a solution to this problem please let us know.

**See Also**

[Forum Topic](<http://www.suneido.com/forum/topic.asp?TOPIC_ID=1215>)