### Print

``` suneido
(@args)
```

Displays the arguments on the active WorkSpace's output pane separated by spaces. Named arguments print the name as well as the value, but their order will vary. For example:

``` suneido
Print(12, 34, a: 56, b: 78)
	=> 12 34 b: 78 a: 56
```

Non-string arguments are converted using
[Display](<Display.md>)

If Print is called from another thread, it prints to [Trace](<Trace.md>)