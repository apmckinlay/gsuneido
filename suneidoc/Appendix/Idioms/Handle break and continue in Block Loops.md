### Handle break and continue in Block Loops

When you write a function or method that loops calling a supplied block, you should support break and continue.

In a block, break and continue throw "block:break" and "block:continue" respectively.

To handle them write:

``` suneido
try
    block(args)
catch (e, "block:")
    if e is "block:break"
        break
    // else block:continue ... so continue
```

This takes advantage of how the catch pattern can be a prefix.

See also: 
[Blocks](<../../Language/Blocks.md>), 
[break](<../../Language/Statements/break.md>), 
[continue](<../../Language/Statements/continue.md>), 
[try catch](<../../Language/Statements/try catch.md>)