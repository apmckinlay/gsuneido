### DiffListControl

``` suneido
(list1, list2, title1 = "List1", title2 = "List2")
```

Compares two lists of lines of text using [Diff](<../../Language/Reference/Diff.md>) and displays the results side by side..  There are buttons for selecting the Next/Previous change. title1 and title2 are used for the headings.

For example:

``` suneido
DiffListControl(#(one, two, four), #(two, three, four, five))
```

would produce:

![](<../../res/difflist.png>)