### PatternControl

``` suneido
(pattern, width = 10, status = "")
```

PatternControl is a
[FieldControl](<FieldControl.md>)
that validates and formats its input according to one or more *patterns*.

Patterns are separated by '|' and consist of:

|  |  | 
| :---- | :---- |
| a | A letter that will be converted to lower case. | 
| A | A letter that will be converted to upper case. | 
| # | A digit. | 
| < | A digit or letter that will be converted to lower case. | 
| > | A digit or letter that will be converted to upper case. | 
| ^c | A literal c, e.g. ^# for a literal '#'. | 
| ... | Any other characters are taken literally. They are automatically inserted if not entered by the user. | 


For example:

``` suneido
Window(#(Pattern '###-####|###-###-####', width: 12,
    status: "Enter a phone number e.g. 249-5050 or 888-558-5050"))
```

PatternControl is the base class for 
[PhoneControl](<PhoneControl.md>)
and
[ZipPostalControl](<ZipPostalControl.md>).