### range.Plus

``` suneido
( range ) => Range
```

where **range** is a Range type.

Returns a new Range which has its *low* set to the smallest *low* values of the current range and the **range** passed in, and its *high* set to the highest *high* values of the current range and the **range** passed in.

For example:

``` suneido
Range(2, 7).Plus(Range(3, 10)) => Range(2, 10)
Range(2, 7).Plus(Range(1, 5)) => Range(1, 7)
```