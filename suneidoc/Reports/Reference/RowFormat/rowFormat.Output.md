#### rowFormat.Output

Allows "_output" format items within Queries using Row for output.

Arguments may be listed in the same order as the columns, 
or may be named with the Field names of the columns.
For example, if the columns were name, age, and sex,
the following lines would be equivalent:

``` suneido
( _output 'Joe', 23, 'male' )
( _output age: 23, name: 'joe' sex: 'male' )
```