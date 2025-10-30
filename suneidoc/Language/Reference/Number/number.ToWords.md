#### number.ToWords

``` suneido
() => string
```

Returns the number converted to a string of words.

For example:

``` suneido
0.ToWords() => "ZERO"
8.ToWords() => "EIGHT"
18.ToWords() => "EIGHTEEN"
20.ToWords() => "TWENTY"
100.ToWords() => "ONE HUNDRED"
750.ToWords() => "SEVEN HUNDRED AND FIFTY"
1000.ToWords() => "ONE THOUSAND"
1525.ToWords() => "ONE THOUSAND FIVE HUNDRED AND TWENTY FIVE"
9999.ToWords() => "NINE THOUSAND NINE HUNDRED AND NINETY NINE"
10000.ToWords() => "TEN THOUSAND"
80750.ToWords() => "EIGHTY THOUSAND SEVEN HUNDRED AND FIFTY"
100000.ToWords() => "ONE HUNDRED THOUSAND"
785694.ToWords() => "SEVEN HUNDRED AND EIGHTY FIVE THOUSAND SIX HUNDRED AND NINETY FOUR"
1000000.ToWords() => "ONE MILLION"
1555555.ToWords() => "ONE MILLION FIVE HUNDRED AND FIFTY FIVE THOUSAND FIVE HUNDRED AND FIFTY FIVE"
```

**Note:** Any fractional portion (e.g. cents) is ignored.

See also:
[number.ToWordsSimple](<number.ToWordsSimple.md>)