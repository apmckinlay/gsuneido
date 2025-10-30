### [Spy](<../Spy.md>) - Methods
`Return(value1, value2, ..., when: block)`
: Specify the return value(s) for the spied function.   
If multiple values are given, they will be returned consecutively   
If multiple Returns are applied, the spy will return the first return value satisfying its **when** condition

`ClearAndReturn(value1, value2, ..., when: block)`
: Resets the return values for the spy object

`Throw(exeption1, exeption2, ..., when: block)`
: Specify the exception to be thrown for the spied function

`CallLogs`
: Return the call argument history of the spied function as a list of args object

`Close`
: Remove the spy from the spied function