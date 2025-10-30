<div style="float:right"><span class="builtin">Builtin</span></div>

#### file.Readline

``` suneido
() => string or false
```

Returns the next line from a file, or false if at the end.

Looks for newline characters ('\n') as line terminators. The newline is **not** included in the returned string.

Trailing carriage returns ('\r') are also removed from the resulting lines. (Returns within lines are not removed or treated as line terminators.)

For example, input of "one\rONE\ntwo\r\nthree\r\r\nfour" would result in lines of "one\rONE", "two", "three", and "four".

Final lines are returned the same, whether or not they have terminating returns or newlines.

**Note:** There is a line length limit of 4000 bytes. If a line exceeds this, only the first 4000 bytes will be returned, although the entire line will be processed.

For example, to copy a file a line at a time:

``` suneido
File("source")
    {|src|
    File("destination", "w")
        {|dst|
        while false isnt line = src.Readline()
            dst.Writeline(line)
        }
    }
```

See also:
[file.Read](<file.Read.md>)