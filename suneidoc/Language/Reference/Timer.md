### Timer

``` suneido
Timer(reps = 1, secs = false) { ... } => number
```

If secs is false, Timer executes the supplied block reps number of times and returns the elapsed time used.

If secs is supplied, Timer executes the supplied block repeatedly for the specified amount of time and returns the number of repetitions.

Timer uses [Date()](<Date/Date.md>)
and [date.MinusSeconds()](<Date/date.MinusSeconds.md>).
Because of their accuracy you should normally time for a few seconds at least.

For example:

``` suneido
Timer(reps: 100)
    {
    Sleep(10)
    }
=> 1.003

Timer(secs: 3)
    {
    Sleep(10)
    }
=> 304
```