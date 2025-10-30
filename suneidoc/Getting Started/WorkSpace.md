## WorkSpace

When you first start up Suneido you'll see the WorkSpace.  This is the
starting point for the Suneido IDE (Integrated Development Environment).
From here you can write and execute Suneido code 
and launch other Suneido tools such as LibraryView and QueryView.
![](<../res/workspace.png>)
One use of the WorkSpace is to execute chunks of code.  Type the following
into the WorkSpace:

``` suneido
123 + 456
```

You can execute this by pressing F9, or by choosing Run from the Edit menu.
The result (579) will show up in the Output pane at the bottom of the
WorkSpace.

You can execute multiple statements as well as single expressions. Try:

``` suneido
for (i = 0; i < 5; ++i)
     Print(i);
```

Print displays its output in the output pane, so you'll see 0 to 4
displayed.  You can clear the output pane by choosing Clear Output on the Edit
menu.

**Note:** Suneido is case sensitive, for example, the above code must use "Print", not "print".

**Note:** If you have more than one thing entered on the WorkSpace then
you need to select (highlight) the portion of code you want to execute before
pressing F9 or choosing Run.  (If there is no specific selection, F9/Run will
automatically select the non-blank lines around the insertion point.)

<span id="FindByExpression"><b>Find by expression</b></span>

You can use Find > Find by expression to search Suneido expressions (like function calls), seamlessly accounting for line breaks, comments and quotes. You can also utilize a signle-letter wildcard to match any expression and efficiently navigate complex code structures. For example, the search pattern "**Fn(a)**" can match any Fn calls with single argument, such as "**Fn(1)**" or "**Fn(1 + 1)**".

Moreover, this feature introduces enhanced flexibility to argument comparision.

-	The order of named arguments is not sensitive. For example, code "**Fn(arg2: 2, arg1: 1)**" matches the pattern "**Fn(arg1: a, arg2: b)**"
-	Named and unnamed argments can be interchanged. For example, code "**Fn(arg: 1)**" matches both "**Fn(a)**" and "**Fn(arg: a)**", but not "**Fn(arg1: a)**". Similarly, "**Fn(1)**" matches "**Fn(arg: a)**" and "**Fn(a)**"
-	Extra named arguments in the code call are ignored. For example, code "**Fn(arg1: 1, arg2: 2)**" matches the pattern "**Fn(arg1: a)**"
-	For additional examples, please refer to **stdlib:AstSearch_Test**