### ScannerSeq

``` suneido
(s, block)
```

A sequence based on [Scanner](<Scanner.md>).

The optional block can act as both filter and map. It is called with the scanner itself. If the block has no return value, that element is excluded from the sequence. For example:

``` suneido
ScannerSeq("x + y",
    {|scan| scan.Type2() is #WHITESPACE ? Nothing() : scan.Text() })
    => #('x', '+', 'y')
```