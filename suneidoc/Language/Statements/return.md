### return
<pre>
<b>return</b> [ <i>expression</i> [ , <i>expression</i> ... ]]
or
<b>return</b> <b>throw</b> <i>expression</i>
or
<b>return</b> <b>@</b> <i>expression</i>
</pre>

A **return** statement causes execution to leave the current function or method,
optionally returning zero, one, or more values.

**Note**: A return from a [block](<../Blocks.md>) returns from the containing function or method, not just from the block.

**Note:** Normally, attempting to use the result of a function that doesn't return anything will cause an error. However, returning the result of a function call is always valid, regardless of what it returns:

``` suneido
return Func()
```

This will pass on the return value, or lack of one.

If execution reaches the end of a function body,
a return will be performed automatically.
If the last statement is an expression,
then the value of that expression is returned.

When returning multiple values, the values can only be received with a [Multiple Assignment](<Multiple Assignment.md>) statement. In any other context a multiple return does not have a value. **Note**: multiple return/assign is not supported on suneido.js

**return throw** will throw an exception if the *caller* does not "use" the return value, unless it is "" or true.
"use" includes assigning to a variable, passing as an argument, etc. If the return value is a string, that is what is thrown. Otherwise, the exception will be "return value not checked".

It is useful where it is important to check the return value. A function could just throw an exception when it fails, but that means if the caller wants to check the result, they're forced to use an awkward try-catch. With return throw, the caller can check the return value if they want, or if they don't care (or forget) then it will throw when it fails and at least it won't be silently ignored.

For example:

``` suneido
f = function() { return throw "failed" }
Print(f())
=> prints "failed"

f = function() { return throw "failed" }
f()
return
=> throws "failed"

f = function() { return throw false }
f()
return
=> throws "return value not checked"
```

If the function has a successful return value that isn't "" or true, e.g. 0 then you can handle it like:

``` suneido
if result is 0
	return result
return throw result
```

**return @expr** (BuiltDate > 2026-09-21)

- The expression must evaluate to an object (or an error is raised).
- The object must not have any named members (or an error is raised).
- If the object has no values, it is like a bare return (no value).
- If the object has one value, it returns that single value.
- If the object has multiple values, it is like a multiple value return.

This is useful for forwarding or adapting the return value of another call.
For example:

``` suneido
@results = Fn()
Print(:results)
return @results
```

The inverse is [Gather Assignment](<Gather Assignment.md>)