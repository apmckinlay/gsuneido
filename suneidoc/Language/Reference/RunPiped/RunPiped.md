<div style="float:right"><span class="builtin">Builtin</span></div>

#### RunPiped

``` suneido
(command) => pipe
(command) {|pipe| ... }
```

Runs an external program with standard input, output, and error redirected to pipes. The block form closes the pipes on exit from the block.

The block form is preferred since it ensures that the pipes are closed.

For example:

``` suneido
RunPiped("dir")
    { |rp|
    while false isnt s = rp.Read()
        Print(s)
    }

RunPiped("ed")
    { |rp|
    rp.Write("a\n")
    rp.Write("hello world\n")
    rp.Write(".\n")
    rp.Write("Q\n")
    while false isnt s = rp.Read()
        Print(s)
    }
```

**Note:** Some programs buffer their input or output (especially when they notice it is redirected). You may need to do CloseWrite before trying to read.

**Note:** If the program path contains spaces, it must be enclosed in double quotes. The remaining arguments must also be properly quoted and escaped.

See also: [RunPipedOutput](<../RunPipedOutput.md>), [ResourceCounts](<../ResourceCounts.md>)