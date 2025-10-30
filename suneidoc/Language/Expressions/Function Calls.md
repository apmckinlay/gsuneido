### Function Calls

function-call:

``` suneido
func(arguments named-arguments)
func(@object)
func(@+# object)
```

argument:

``` suneido
expression
```

named-argument:

``` suneido
name : expression
```

Scalar values such as boolean, number, string, and date are immutable so you can think of them as being passed by value (even though internally they may not be). Mutable values such as objects or class instances are passed by reference (pointer) *so the called code can modify them*.

`@object` passes all the members of the object as arguments. `@+1` object skips over the first member of the object and passes the remainder as arguments.

Function call arguments are matched with function definition parameters as follows:

-	if both argument and parameter use @ then a copy of the object is passed directly
-	else if @parameter then the arguments are used to construct an object
-	else if @argument this is treated as if the members of the object were passed as arguments
-	unnamed arguments are assigned in order to the parameters, 
	then named arguments are assigned to parameters by name 
	(potentially, a named argument could replace the value assigned to a parameter by an unnamed argument)
-	finally, default values are assigned to uninitialized parameters
-	if any parameters do not have values after this process, a "missing argument" error will occur


Commas are optional between arguments
except where they are necessary to separate ambiguous sequences.
For example:

``` suneido
fn(a + b)
fn(a, +b)
```

**Note:** A [block](<../Blocks.md>) immediately following a function call is interpreted as additional argument named "block". For example:

``` suneido
fn(x)
    { ... }
```

is equivalent to:

``` suneido
fn(x, block: { ... })
```
**Inverted Method Calls**`string(value ...)` is treated as `value.name(...)`For example, "Size"("hello") is treated as "hello".Size()
This is primarily useful for functional style programming, for example:

``` suneido
ob = #(one, two, three, four)
ob.Map(#Size) // or ob.Map("Size")
    => #(3, 3, 5, 4)
```

There is a shortcut for passing a variable as a named argument:
`func(:a, :b)` is equivalent to `func(a: a, b: b)` This is useful when printing some variables e.g. `Print(:i, :j)`, when constructing objects or records e.g. `Object(:name, :age)` or for super calls e.g. `super(text, :width, :size)`

See also [Functions](<../Functions.md>)