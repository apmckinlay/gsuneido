#### report.Construct

``` suneido
( format_spec ) => format
```

Creates an instance of a format object.

Normally accessed as _report.Construct.

Takes a string containing the name of a format, 
or an object whose first member is the name of a format.
If the name starts with a lower case letter, 
it is assumed to be a field name and the format is looked and merged with the supplied format.
Otherwise, "Format" is appended onto the name.
Then Construct is used to create the instance.

Automatically copies the following members to the new instance with their names capitalized:

``` suneido
x, y, xmin, ymin, xstretch, ystretch, field, span
```

Note: x, y, xmin, ymin values are assumed to be in inches 
and are multiplied by 1440 to convert twips.