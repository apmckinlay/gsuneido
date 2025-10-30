### LibraryStrings

``` suneido
(library, table = "translatelanguage") => number
```

Prints a list of strings found in the library that are not in the translation table.

Ignores:
    
strings less than 5 characters long
 
strings with characters other than letters, periods, and ampersands
 
strings that start with www., create, update, alter, delete, insert, or destroy

It removes trailing '...', &'s, and leading and trailing whitespace.

Returns the number of strings printed.