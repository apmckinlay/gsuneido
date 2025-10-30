### Stopwatch

``` suneido
() => stopwatch

stopwatch() => string
```

Records the time when it was created and reports the elapsed and time using [ReadableDuration](<ReadableDuration.md>).

For example:

``` suneido
sw = Stopwatch()
Sleep(500)
sw()

=>  510 ms

sw = Stopwatch()
Sleep(500)
Print(sw())
Sleep(2000)
Print(sw())
=>  510 ms
    + 2.01 sec = 2.52 sec
```