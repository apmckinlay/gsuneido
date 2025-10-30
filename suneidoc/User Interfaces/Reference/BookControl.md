### BookControl

``` suneido
(book = false, title = false, start = "Cover", login = false,
    help_book = false)
```

Provides an interface to navigate through Suneido programs and html pages which are stored in a book table (see the [Getting Started](<../../Getting Started.md>) section for help on creating a book).

The **book** argument is the name of the book table to be used.

**title** is what will be displayed in the book window title bar.

The **start** argument should be the name of one of the book pages at the root level.  This defaults to "Cover", so if you define a page called "Cover" at the root level of the book then this will be the first page displayed when the book is started.

The **login** argument (a function or a string that evaluates to a function) is used to specify code to call that will handle any login procedures that need to be done.  If the login function returns false, the book window will be destroyed.

If **help_book** is true, then Find and Print buttons will be shown.

If there is help for the specified book (named: book $ "Help") then a Help button will be displayed.

For example, the following will create a BookControl with "Suneido Documentation" in the window's title bar, starting on the Introduction page.

``` suneido
BookControl('mybook', start: "Contents")
```

Would display:
![](<../../res/mybook.png>)
If a page does not begin with \<html> HtmlPrefix and HtmlSuffix will be added to the beginning and end respectively. The default definitions in stdlib are:

|  |  |  |  |  | 
| :---- | :---- | :---- | :---- | :---- |
| `HtmlPrefix` | <code>&lt;html><br />&lt;head><br />&lt;/head><br />&lt;body></code> |  | `HtmlSuffix` | <code>&lt;/body><br />&lt;/html></code> | 


You can provide your own definitions by creating HtmlPrefix and/or HtmlSuffix pages in the res folder of your book. For example, suneidoc (this User's Manual) uses HtmlPrefix to define a style sheet.

**Note:** If *book*.ico (where *book* is the name of the book) exists in the same directory as [ExePath](<../../Language/Reference/ExePath.md>), then it will be loaded and set as the window icon.

See also:
[Books](<../../Tools/Books.md>)