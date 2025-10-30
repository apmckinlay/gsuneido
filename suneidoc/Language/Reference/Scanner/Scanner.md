<div style="float:right"><span class="builtin">Builtin</span></div>

#### Scanner

``` suneido
(string) => scanner
```

Returns a scanner object on the string that tokenizes the same way as the Suneido language compiler.

For example, to remove comments:

``` suneido
text_in = "hello /* there */ world"
text_out = ""
scan = Scanner(text_in)
for token in scan
    if scan.Type2() isnt #COMMENT
        text_out $= scan.Text()
return text_out

=> "hello  world"
```