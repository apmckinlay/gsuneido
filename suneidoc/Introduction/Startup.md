## Startup

When Suneido starts up, it first tries to open the database (suneido.db) in the current working directory.

If suneido.db is **not** found, it tries to run the command line, either as the name of a file containing a function or class, or as code to be evaluated. As with VersionMismatch, the code can only use built=in capabilities - it has no database or libraries.

If the database is opened successfully, it automatically [Use](<../Language/Reference/Use.md>)'s "stdlib". It then calls Init (or Init.Repl for the REPL), which must be defined in stdlib.

The standard Init does the following:

If there are no command line arguments and if not a server, the default (IDE) persistent set is loaded.

If there are any command line arguments, Init tries to interpret them several different ways:

-	as the name of a persistent set to load   
	e.g. `suneido IDE`
-	as the name of a text file containing code to execute   
	e.g. `suneido server.go`
-	as code to be executed   
	e.g. `suneido Login("IDE")`

On Windows there are two versions of Suneido:
gsuneido.exe
: This is a 
[Windows GUI program](<https://stackoverflow.com/questions/574911/difference-between-windows-and-console-application>). It includes the Win32 interface. It automatically runs the Windows message loop. stdout and stderr are appended to error.log

gsport.exe
: This version should be used for the server. It does not contain the Win32 interface.

On other operating systems e.g. Linux or Mac, the executable is similar to gsport.exe.

When running client-server, the built date must match. If the dates do not match, the server will look for a record in stdlib called `VersionMismatch`. If it finds it, it will send it to the client, that will then execute it. This can be used, for example, to handle automatically updating the client. **Note:** the VersionMismatch code will run without a database so it can only use built-in functions and classes.

When GUI mode (gsuneido.exe) or [Running as a Service](<Running as a Service.md>), error output is redirected to a file. By default this is "error.log" in the current directory, but if running as a client on Windows it will be \<appdata>/suneido\<port>.err and on other systems it will be \<tempdir>/suneido\<port>.err