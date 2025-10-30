### TwoListDlgControl

``` suneido
(list, initial_list, title)
```

Similar to [TwoListControl](<TwoListControl.md>) but it it has OK and Cancel buttons on the bottom right corner. 

For example:

``` suneido
Window(#(TwoListDlg (1 2 3) (4 5 6 7) 'TwoListDlg'))
```

would display

![](<../../res/TwoListDlg.gif>)