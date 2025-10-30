## Differences from JavaScript

In Suneido:

-	You normally use .name rather than this.name
-	Only "blocks" are closures. functions do not have access to their enclosing lexical context.
-	`return` from a block returns from the defining function, not just from the block.
-	Setting a member on an instance of a class will set it in the instance, even if the member exists on a super class.
-	Classes are read-only.
-	The only inheritance is with classes and instances. There is no equivalent to JavaScript prototypes.
-	Code is in library tables in the database. Each record in a library defines a single immutable constant e.g. number, string, object, function, class.
-	Suneido Objects combine an array (unnamed values) and a map (named values). Named values are similar to JavaScript object properties. The array does not have "gaps". Any values after a gap will be in the named members.
-	Only Objects/Records and class instances have named values. There is no equivalent to JavaScript properties on other types of values.
-	There are no "primitive" values in Suneido.
-	Suneido numbers are *decimal* floating point. Numbers like .1 are represented exactly, without floating point approximation.
-	No null or undefined. Internally Suneido has "uninitialized" values, but any attempt to use them will throw an exception so they are never "visible".