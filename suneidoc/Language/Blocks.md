## Blocks

A block is a section of code within a function. It can be called like a function, and can have parameters and accept arguments like a function. But blocks can share local variables with the containing function like a closure.

Blocks can be used to create user defined control structures or as anonymous functions.

Blocks are values. This means they can be assigned to variables, passed to other functions, etc.

Blocks are written as:

``` suneido
{ ... }
{|x, y| ... }
{|@args| ... }
```

Commas between parameters are optional. Blocks have the same parameters as functions and methods except they do not currently allow default values.

Block parameters are independent from local variables with the same names in the containing function. i.e. If you have a block parameter called "x", the block will not be able to access "x" in the containing function, and any changes the block makes to its "x" will not affect "x" outside the block. These block parameters will appear in the debugger preceded by an underscore.

The "value" of a block is the value of its last statement. On the other hand, an actual "return" will return from the function containing the block.

For example:

``` suneido
for_each = function (ob, block)
    {
    for (x in ob)
        block(x)
    }

sum = 0
for_each(#(1, 2, 3, 4), {|x| sum += x; })
Print(sum)
    => 10
```

A block immediately following a function call is interpreted as another argument. So the above example could also be written as:

``` suneido
for_each(#(1, 2, 3, 4))
    {|x| sum += x; };
```

A block following the argument parenthesis will be passed as `block: block`. This means the parameter that receives it <u>must</u> be called "block".

Empty argument parenthesis may be omitted in most cases, for example these are all equivalent:

``` suneido
10.Times({ F() })
10.Times() { F() }
10.Times { F() }
```

Another shortcut is that if a block has no parameters and it refers to a variable called "it" or "_" then a parameter will be automatically added to the block. For example:

``` suneido
#(12, 34, 56).Each { Print(it) }
=>  12
    34
    56

#(1, 2, 3).Map { _ * 2 }
=>  #(2, 4, 6)
```

The exception is that `Name { ... }` is treated as `class : Name { ... }` rather than `Name({ ... })`

One of the powerful aspects of blocks is that they can outlive the function call that created them. For example:

``` suneido
make_counter = function (next)
    { return { next++ } }
counter = make_counter(10)
Print(counter())
Print(counter())
Print(counter())
    =>  10
        11
        12
```

Within a block, `break` does `throw "block:break"` and `continue` does `throw "block:continue"`. This allows user defined looping control structures to handle them. For example, here is the transaction.QueryApply method:

``` suneido
QueryApply(query, block, dir = 'Next')
    {
    .Query(query)
        { |q|
        while (false isnt x = q[dir]())
            try
                block(x)
            catch (ex, "block:")
                if ex is "block:break"
                    break
                // else block:continue ... so continue
        }
    }
```

**Note:** If block:break and block:continue exceptions are not caught and break or continue are used you will get an error.

While a block looks similar to a function, blocks have to be constructed. For example:

``` suneido
b = { ... }
```

is actually implemented more like:

``` suneido
b = make_block(...)
```

i.e. an instance of a block must be constructed at runtime to associate it with the current environment (parent local variables).

Whereas:

``` suneido
f = function () { ... }
```

does not need to create an instance, there is just one function value, created at compile time.

Note: Suneido optimize blocks that do not reference parent local variables to be just regular functions, not closures.

**Warning:** if a closure becomes concurrent (by being stored in a concurrent object) it becomes "detached" from its original context and any changes will not affect its origin. For example:

``` suneido
x = 1
b = { x++ }
b() // modifies x
Print(x)
	=> 2
Suneido.b = b // this makes the block concurrent
b() // does not modify the original x
Print(x)
	=> 2
```