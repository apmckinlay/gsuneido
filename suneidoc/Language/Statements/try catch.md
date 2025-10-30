### try catch
<pre>
<b>try</b>
    <i>statement1</i>
[ <b>catch</b> [ ( variable [ , pattern ] ) ]
    <i>statement2</i> ]
</pre>

Executes the **try** *statement1*. If any exceptions occur, then the **catch** *statement2* is executed. If there is no catch, exceptions during *statement1* are ignored.

If a **variable** name is specified in the catch it will be set to the exception string. For example:

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

For example, to ignore errors if the table doesn't exist:

``` suneido
try
    Database("drop mytable")
```

or to handle a value not being numeric:

``` suneido
try
    x = Number(x)
catch
    x = 0
```

See also:
[throw](<throw.md>),
[Exceptions](<../Exceptions.md>)