<div style="float:right"><span class="builtin">Builtin</span></div>

#### File

``` suneido
(filename, mode = "r") => file
(filename, mode = "r") {|file| ... }
```

Create a file object for accessing a file.

File is based on the standard stream io functions 
(e.g. fopen, fclose, fflush, fread, fwrite, fseek, ftell).

The *mode* specifies the type of access for the file:
`**"r"**`
: Open a file for reading. 
If the file does not exist or cannot be opened, an exception will be thrown:
"File: can't open '*filename*' for r".

`**"w"**`
: Create an empty file for writing. 
If the file exists, its contents are destroyed.

`**"a"**`
: Open a file for writing at the end of the file (appending); 
creates the file if it doesn’t exist.

When a file is opened with `"a"`, 
all write operations occur at the end of the file.
This means existing data cannot be overwritten.
Seek is not allowed.

It's normally safer to use the block form of File 
to ensure the file gets closed even if there are exceptions.

For example, to process a file 1000 bytes at a time:

``` suneido
File(source)
	{|src|
	File(destination, "w")
		{|dst|
		while false isnt s = src.Read(1000)
			{
			s2 = Process(s)
			dst.Write(s2)
			}
		}
	}
```