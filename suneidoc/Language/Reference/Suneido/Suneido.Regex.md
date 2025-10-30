<div style="float:right"><span class="builtin">Builtin</span></div>

#### Suneido.Regex

``` suneido
(string) => regex
```

Returns a compiled regular expression.

For debugging purposes, a compiled regular expression can be disassembled. For example:

``` suneido
Suneido.Regex("^a?b").Disasm()
=>	0: LineStart
	1: SplitNext 6
	4: Char a
	6: Char b
	8: DoneSave1
```

**NOTE**: This should never be needed since compiled regular expressions are cached.