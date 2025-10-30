<div style="float:right"><span class="builtin">Builtin</span></div>

### Spawn

``` suneido
(mode, command, @arguments) => number
```

Create and execute a new process. 
mode
: 

|  |  | 
| :---- | :---- |
| P.WAIT | Suspends calling thread until execution of new process is complete. | 
| P.NOWAIT | Continues to execute calling process concurrently with new process. | 

command
: Specifies the file to execute.

arguments
: Zero or more string arguments to be passed to the command. Separate command line arguments should be separate function arguments, do not combine multiple command line arguments into one string.

For example:

``` suneido
Spawn(P.NOWAIT, "notepad", "tmp.txt")
```

The return value from a synchronous Spawn (P.WAIT) is the exit status of the process (0 usually means the process terminated normally). The return value from an asynchronous Spawn (P.NOWAIT) is the process id (or -1 if it fails to start the process).

For P.WAIT, if the return value is not 0 (zero) it is [return-throw](<../Statements/return.md>), i.e. if the result is not used it will throw an exception. (as of BuiltDate 2025-01-16)

**Note:** On Windows, if you use Spawn to run a console program (rather than a gui program) it will open a console window. If you don't want this, use [RunPiped](<RunPiped.md>)

See also:
[System](<System.md>), 
[RunPiped](<RunPiped.md>)