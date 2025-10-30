<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Compile

``` suneido
(object = false) => value
```

Compiles the string, which must be a constant (e.g. number, function, class). Preferable to [string.Eval](<string.Eval.md>) because it does not execute arbitrary code.

For example:

``` suneido
"function () { }".Compile()
    => /* function */
```

If an object is supplied as an argument then the offsets of "warnings" (e.g. uninitialized variables) are added to it. This is used by CheckCode which is used by [LibraryView](<../../../Tools/LibraryView.md>) to show mistakes.

**Note**: Compile may be slower when asking for warnings (by supplying an object argument) because it may load library records as part of determining whether global names are defined.

Errors and warnings will be added to the object as positions in the source string. Warnings will be negative, errors will be positive. Currently, the following are detected:

-	local variable used but not initialized (error)
-	local variable initialized (including as function parameter) but not used (warning)
-	reference to an undefined global name (error)
-	_Name where Name is defined (warning - may be invalid in context)   
	(_Name where Name is undefined will throw an exception from Compile)


Only the starting position is given. To get the length or the text, you can use [Scanner](<../Scanner.md>)

``` suneido
scanner = Scanner(source[pos ..])
scanner.Next()
text = scanner.Text()
```

See also:
[string.Eval](<string.Eval.md>)