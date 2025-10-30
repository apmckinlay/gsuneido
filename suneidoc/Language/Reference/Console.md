### Console

A simple interface to the Windows console functions.  It can be used like:

``` suneido
con = Console()
con.SetTitle("Suneido Console")
con.Write("What is your name? ")
name = con.Readline()
con.Writeline("Hello " $ name)
con.Writeline("Press ENTER to quit")
con.Readline()
con.Close()
```

Has the following methods:
`SetTitle(string)`
: Uses SetConsoleTitle to set the name on the titlebar of the console.

`Write(string)`
: Output the string to the console.

`Writeline(string)`
: Output the string to the console, followed by carriage return, newline (\r\n).

`Readline(string)`
: Read a line of input from the console. Readline will not return until the user presses Enter. All other Suneido activity will be blocked during this time. Note: This assumes the console is in line oriented input mode.

`Close()`
: Closes the console window.