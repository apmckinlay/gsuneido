## Libraries

Libraries contain definitions of user defined global names. A library definition must be a literal (e.g. number, string, object, function, class).

Suneido automatically uses the standard library "stdlib".

Within the redefinition of a function or class in a library, the previous definition of the global name may be referred to by preceding it with an underscore.

**Note:** _Name will only work for the Name you are re-defining. i.e. Within a re-definition of X you can refer to _X but you cannot refer to _Y

**Warning:** Currently, you cannot save (e.g. store in the database) compiled library records that refer to _Name's.

The first time a global name is referenced, each of the libraries in use will be searched for its definition.

If a library record is added, modified, or deleted, that global name should be unloaded (so it will be re-loaded on demand). LibraryView handles this automatically, but if you modify a library record through other means (e.g. QueryView) then you should Unload it yourself.

Note: Built-in global names can <u>not</u> be redefined in libraries.

See also:
[Use](<Reference/Use.md>),
[Libraries](<Reference/Libraries.md>),
[Suneido.LibraryTags](<Reference/Suneido/Suneido.LibraryTags.md>),
[Unload](<Reference/Unload.md>),
[Unuse](<Reference/Unuse.md>)