## Exceptions
<pre>
<b>try</b>
    <i>statement1</i>
[ <b>catch</b> [ ( variable [ , pattern ] ) ]
    <i>statement2</i> ]

<b>throw</b> <i>string</i>
</pre>

Within the **try** statement and any functions called by it (directly or indirectly), a throw will cause an immediate transfer of control to the most recent **catch** statement whose prefix matches the exception.  If no prefix is specified, the catch will match all exceptions.

**NOTE**: try-catch will not work properly if nested within the same function.

**throw** may be used within the **catch** statement to pass the exception to the next catch.

Note: Suneido also generates exceptions internally for things like uninitialized variables and syntax errors.

An uncaught exception will call a user defined Handler function, passing it the call stack as a list of objects, each containing fn (as returned by [Frame](<Reference/Frame.md>)) and locals (as returned by [Locals](<Reference/Locals.md>)). The Handler in stdlib will open a DebuggerControl showing the function calls leading up to the exception if Suneido.User is 'default', otherwise it will log the error with SuneidoLog and call AlertError to display "Program Error" to the user.

If a **variable** name is specified in the catch it will be set to the exception. For example:

``` suneido
try
    throw "oops"
catch (e)
    Print(caught: e)

=> caught: oops
```

A **catch** with a **pattern** string will catch any exception that starts with that pattern (has it as a prefix). If the pattern starts with a "*" it will match anywhere in the exception. Multiple patterns may be specified by separating them with "|". For example:

``` suneido
try
    MyFunc()
catch (e, "*error|*failure")
    ...
```

This would catch exceptions that contained "error" or "failure". Any other exceptions would not be caught.

Note: The catch pattern must be a literal string - it cannot be a variable or expression.

The value passed to **catch** is a special exception type. It can be used as a normal string, but it also has two extra methods:
exception.Callstack()
: Returns the call stack (as a list of objects, with fn the function/block/method and locals, an object with members for each local variable) at the point of the exception. This is useful when logging exceptions. For example:

``` suneido
try
    MyFunc()
catch (e)
    Log(e, e.Callstack())
```
: You can also use DebuggerControl to view the callstack information since it is in the same format as the call stack passed to Handler. For example:

``` suneido
DebuggerControl(e, e.Callstack(), 0)
```

exception.As(string) <span class="deprecated">Deprecated</span>
: Returns a new exception value with the specified string, but keeping the call stack information. This is useful when catching and re-throwing an exception. For example:

``` suneido
try
    MyFunc()
catch (e)
    throw e.As("error in MyFunc: " $ e)
```

**Note**: Concatenation preserves exception information eliminating the need for most uses of As.

Note: Within [blocks](<Blocks.md>), break and continue throw "block:break" and "block:continue". If you are writing a loop type function that accepts a block, you can implement break and continue by catching and handling these exceptions. For example:

``` suneido
while ...
    try
        block(x)
    catch (e, "block:")
        if e is "block:break"
            break
        // else block:continue ... so continue
```

**WARNING**: Currently Suneido does <u>not</u> handle nested try-catch within a single function/method. i.e. the following example will <u>not</u> catch the exception.

``` suneido
try
    try
        throw "should be caught"
    catch (unused, 'something else')
        {}
catch
    {}
```

It is ok to nest separate functions/methods that each have try-catch so one workaround is to extract the inner try-catch into its own function/method.