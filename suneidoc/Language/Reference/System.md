<div style="float:right"><span class="builtin">Builtin</span></div>

### System

``` suneido
(command) => number
```

Passes *command* to the command interpreter, which executes the string as an operating-system command. Uses the COMSPEC or SHELL environment variables 
to locate the command-interpreter.

System returns the value that is returned by the command interpreter. Normally a return value of 0 (zero) means success, and -1 indicates an error.

If the return value is not 0 (zero) it is [return-throw](<../Statements/return.md>), i.e. if the result is not used it will throw an exception. (as of BuiltDate 2025-01-16)

For example:

``` suneido
System("notepad")
```

See also:
[Spawn](<Spawn.md>), 
[RunPiped](<RunPiped.md>)