### Extracting Text Lines

``` suneido
i = s.Find1of("\r\n")
line = s[..i]
if (s[i..2] is "\r\n")
    ++i
s = s[i + 1..]
```

This code extracts the next line of text from s, accepting line endings of "\r", "\n", or "\r\n". The line ending is not included on the line.