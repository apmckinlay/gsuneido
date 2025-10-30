### RunPipedOutput

``` suneido
(command, input = "") => string
```

Runs command, writes input to it, and returns the output (both stdout and stderr).

To get the exit value, use:

``` suneido
RunPipedOutput.WithExitValue(command) => Object(:output, :exitValue)
```

WithExitValue also takes an optional input argument.

See also: [RunPiped](<RunPiped.md>)