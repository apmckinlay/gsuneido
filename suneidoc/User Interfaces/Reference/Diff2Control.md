### Diff2Control

``` suneido
CallClass(title, srcOld, srcNew, titleOld, titleNew)

New(srcOld, srcNew, lib, recName, titleNew = "ListNew", titleOld = "ListOld)
```

Compares two sources (text or lists of lines) using [Diff](<../../Language/Reference/Diff.md>) and displays the results side by side..  There are buttons for selecting the Next/Previous change. title1 and title2 are used for the headings.

For example:

``` suneido
Diff2Control("Diff", #(one, two, four), #(two, three, four, five), "List1", "List2")
```

would produce:

![](<../../res/difflist.png>)

CallClass creates a Window. You can use Dialog (and New):

``` suneido
Dialog(0, ["Diff2", "hello\nworld", "hello\nthere", "", ""])
```