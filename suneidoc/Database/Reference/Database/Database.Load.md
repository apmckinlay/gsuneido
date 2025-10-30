<div style="float:right"><span class="builtin">Builtin</span></div>

### Database.Load

``` suneido
(tablename, fromfile = "", privateKey = "", passphrase = "") => number
```

Loads a table into the current database. If **fromfile** is "" (or omitted) it defaults to "table.su". If client-server, this happens on the server.

If **privateKey** and **passphrase** are supplied, the file is decrypted as it is loaded, equivalent to [OpenPGP](<../../../Language/Reference/OpenPGP.md>).PublicDecrypt.

Equivalent to the `-load` [command line option](<../../../Introduction/Command Line Options.md>)

Dump files are created with the `-dump` [command line option](<../../../Introduction/Command Line Options.md>) or with [Database.Dump](<Database.Dump.md>)

**Warning**: If you load a library table that is currently in use you should call [Unload](<../../../Language/Reference/Unload.md>)() after loading.
Note: Unlike Database.Dump, Database.Load cannot be used to load an entire database. Use the `-load`
[command line option](<../../../Introduction/Command Line Options.md>).