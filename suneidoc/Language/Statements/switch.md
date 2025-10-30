### switch
<pre>
<b>switch</b> [ <i>expression</i> ]
    {
<b>case</b> <i>literals</i> :
    <i>statements</i>
...
[ <b>default :</b>
    <i>statements</i> ]
    }
</pre>

Executes the first matching case, or if no cases match, then the default, if there is one. A case matches if the value of the expression is equal to one of the literal values listed for the case.

For example:

``` suneido
switch i
    {
case 0, 1, 2 :
    Print("i is 0, 1, or 2")
case 3 :
    Print("i is 3")
default :
    Print("invalid value of i")
    }
```

**Note:** Suneido's switch statement differs from C, C++, or Java in that it allows multiple values on a single case, and break statements are not used to prevent execution from *falling through* to the next case. It also supports more types of value, e.g. including strings and dates.

**Note:** If there is no default and the expression does not match any of the cases, then an exception of "unhandled switch value" will be thrown. This can be disabled by setting Suneido.SwitchUnhandledThrow = false (see [Compiler Options](<../../Introduction/Compiler Options.md>))

The switch expression is optional, it defaults to true. In this case switch becomes equivalent to a series of if-else but is sometime clearer written this way.

``` suneido
switch
    {
case i < 0 :
    ...
case i is 0 :
    ...
case i > 0 :
    ...
    }
```