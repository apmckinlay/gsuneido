<div style="float:right"><span class="builtin">Builtin</span></div>

#### file.Close

Close a file object.

File's should always be closed when you have finished using them.

**Note:** You cannot rename or delete a file while it is open.

For example:

``` suneido
f = File("tmp", "w")
f.Writeline("hello world")
f.Close()
```

**Note:** It's often safer to use the block form of File 
to ensure the file gets closed even if there are exceptions.
For example:

``` suneido
File("tmp", "w")
    { |f|
    f.Writeline("hello world")
    }
```