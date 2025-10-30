### Working

``` suneido
(message, block)
(message) { ... }
```

Working displays a window containing the specified message
while the block is executed.

For example:

``` suneido
Working("sleeping...")
	{
	Sleep(2000)
	}
```