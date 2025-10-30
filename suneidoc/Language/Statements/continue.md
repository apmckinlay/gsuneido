### continue

The **continue** statement can be used to end the current iteration and go on to the next iteration of a loop (while, do-while, for, forever).

For example, to skip processing values less than zero:

``` suneido
for x in ob
    {
    if x < 0
        continue
    // process positive values
    ...
    }
```

See also:
[break](<break.md>)