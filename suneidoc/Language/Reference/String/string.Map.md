#### string.Map

``` suneido
(block) => string
```

Calls block(c) for each character in the string and returns the results concatenated together.

For example:

``` suneido
"Hello World".Map({|c| c is 'l' ? 'L' : c})
    => "HeLLo WorLd"
```