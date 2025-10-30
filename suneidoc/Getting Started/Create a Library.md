## Create a Library

If you're going to write any code, the first thing you need to do is create
your own library.  Choose New Library from the LibraryView File menu and enter
`mylib` as the name of your library.  The usual convention is to
make table names valid identifiers, and for library tables to end with "lib". 
New Library will create the library and 
[Use](<../Language/Reference/Use.md>) it.

**Note:** If you loaded in the supplied mylib, just choose Use Library
from the File menu and select `mylib` from the list.

The order that libraries are used is significant.  When Suneido is looking
up the definition of a global name, the libraries are searched starting with
the most recently used and ending with stdlib, which is always used when
Suneido starts up.  One result of this is that later libraries can redefine
things in earlier libraries.  This *layering *of libraries is an
effective method of tailoring software by having a base library and then a
customization library that redefines or extends the base library.

**Note: **If you want to modify stdlib records, you should copy the
records to your own library and modify them there.  That way, if you get a new
version of stdlib, you won't lose your changes.  And if your changes don't
work, you can always just delete the record from your library. ** Tip:**
You can drag items in the tree, or right click on them to copy and paste.  To
select multiple items hold down the Ctrl key while clicking.

Suneido's persistence mechanism saves which libraries are in use so if you
exit and then restart Suneido, both your library (mylib) and stdlib will be in
use.  To stop using a library, choose Unuse Library from the Library View File menu and select the library from the list.