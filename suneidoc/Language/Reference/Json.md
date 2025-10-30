### Json

Encodes and decodes between Suneido values and JSON strings.

``` suneido
Json.Encode(value)
    => json_string

Json.Decode(json_string)
    => value
```

Suneido does not have separate data types for arrays and maps, and allows you to mix both. An object with only unnamed list members will be encoded as a JSON array [...]. An object with named members will be encoded as a JSON object {...}.

`Json.Decode(Json.Encode(value))` may not return the same value. For example:

``` suneido
Json.Decode(Json.Encode(#20171218)) => "#20171218"
```

In this case the Suneido date is converted to a string by Encode, but JSON strings are not converted back to dates.

Decode may throw "Invalid Json format"

See Json_Test in stdlib for examples.