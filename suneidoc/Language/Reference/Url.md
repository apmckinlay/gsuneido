### Url

Utility methods for dealing with URL's.
`BuildQueryString(object) => string`
: The reverse of ResolveQueryString. Joins values with '&', named members to '=', encodes values with EncodeQueryValue.

`Decode(string) => string`
: Converts '+' to space and decodes %xx hex escapes.

`Encode(url, object = #()) => string`
: Converts space to '+'. Converts characters other than +#&*-./:;=?@_0-9a-zA-Z to %xx hex escapes.

`EncodeQueryValue(string)`
: Converts characters other than -_.~0-9a-zA-Z to %xx hex escapes.

`ResolveQueryString(string) => object`
: The reverse of BuildQueryString. Splits on '&', converts '=' to named members, applies Decode to the values, converts numeric strings to number.

`Split(url) => object`
: Splits the url into scheme, host, port, user, fragment, path, basepath, and query.

See also: [Paths](<Paths.md>)