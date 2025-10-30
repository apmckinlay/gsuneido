### Type

``` suneido
Type(value) => string
```

Returns the type of the value. For example:

``` suneido
Type(true) => "Boolean"
Type(123) => "Number"
```

Types include:

-	Boolean
-	Number
-	String
-	Date
-	Object (includes class instances)
-	Record
-	Function
-	Block
-	Class
-	Method (result of retrieving method from class)
-	Transaction
-	Query
-	Cursor
-	COMobject
-	Builtin
-	BuiltinClass


See also:
[Boolean?](<Boolean?.md>),
[Class?](<Class?.md>),
[Date?](<Date?.md>),
[Function?](<Function?.md>),
[Number?](<Number?.md>),
[String?](<String?.md>),
[Object?](<Object?.md>),
[Record?](<../../Database/Reference/Record?.md>)