<div style="float:right"><span class="builtin">Builtin</span></div>

#### number.Format

``` suneido
(format) => string
```

Returns the number converted to a formatted string.  If the mask is not long enough to handle the number (i.e. not enough digits) then "#" will be returned.  

|  |  | 
| :---- | :---- |
| `#` | display a digit, zero fill if on the right of decimal | 
| `-` | displayed if number is negative, otherwise converted to space | 
| `( )` | displayed if number is negative, otherwise converted to space | 
| `,` | displayed if within number | 
| `.` | displayed if within number | 


For example:

``` suneido
n = 1234;
m = -1234;
n.Format("##")          => '#'              // doesn't fit
n.Format("####")        =>  '1234'
n.Format("###,###")     =>  '1,345'
m.Format("####")        =>  '-'             // format doesn't handle negatives
n.Format("(####)")      =>  '1234 '
m.Format("(####)")      =>  '(1234)'
d = 123.456;
d.Format("###.###")     =>  '123.456'
d.Format("###")         => '123'
d.Format("###.##")      => '123.46'         // rounded
d.Format("###.#####")   => '123.45600'
```

Note, the resulting string is not left padded with leading blanks to make it a fixed length - it should be displayed right justified.