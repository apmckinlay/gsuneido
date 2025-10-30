<div style="float:right"><span class="builtin">Builtin</span></div>

#### socketClient.Readline

Returns the next line from the socket connection.

``` suneido
() => string
```

Returns the next line from the socket connection.

Looks for newline characters ('\n') as line terminators. The newline is **not** included in the returned string.

Trailing carriage returns ('\r') are also removed from the resulting lines. (Returns within lines are not removed or treated as line terminators.)

For example, input of "one\rONE\ntwo\r\nthree\r\r\nfour" would result in lines of "one\rONE", "two", "three", and "four".

Final lines are returned the same, whether or not they have terminating returns or newlines.

**Note:** There is a line length limit of 4000 bytes. If a line exceeds this, only the first 4000 bytes will be returned, although the entire line will be processed.

In case of a lost connection or timeout, if partial data has been received it will be returned. i.e. the returned string may not always be a complete line.

An exception will be thrown (`"lost connection or timeout"`) if the connection is lost or the timeout is exceeded.

**Note:** There is a maximum line length of 2000 characters. If a line exceeds this length the first 2000 characters will be returned and the remainder of the line will be discarded.