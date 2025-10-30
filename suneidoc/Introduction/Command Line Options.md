## Command Line Options

If the command line is empty, Suneido will look for a file named "suneido.args", first in the current directory, and if that fails, in the executable directory. The suneido.args file should contain a single line with the command line arguments.

`-check`
: Verify the integrity of the database.

`-c[lient] address`
: Run Suneido as a client, connecting to the specified IP address.  If -c (or -client) is the last option or is followed by -- or another option (e.g. -c -p 3150) the address defaults to 127.0.0.1 - the local machine.

`-compact`
: Remove unused space (e.g. deleted information) from the database.

`-d[ump]`
: Dump the entire database to database.su

`-d[ump] tablename`
: Dump the specified table to tablename.su e.g. `-dump stdlib` would dump stdlib to stdlib.su.   
See also: 
[Database.Dump](<../Database/Reference/Database/Database.Dump.md>)

`-e[rr]p[ort]=port`
: Set the error log path as if running as a client with the specified port. i.e. on Windows <appdata>/suneido<port>.err and on other systems <tempdir>/suneido<port>.err     
**Note**: The error log path is only used when gui mode (gsuneido.exe) or when running as a service.

`-h[elp] -?`
: Display a list of the valid options.

`-l[oad] [@filename]`
: Load the entire database from database.su (or the specified **filename**) and renames the old database to suneido.bak

`-l[oad] tablename`
: Load the specified table from tablename.su e.g. `-load stdlib` would load stdlib from stdlib.su   
See also: 
[Database.Load](<../Database/Reference/Database/Database.Load.md>)
: **Warning:** Destroys any existing table with the specified name.

`-p[ass]p[hrase]=passphrase`
: Only used with **-load**. A private key is read from stdin and the file is decrypted as it is loaded.

`-p[ort]=#`
: Choose a specific TCP/IP port for client or server. The default port is 3147. The web server monitor port is one higher than the main server port, i.e. the default is 3148.

`-repair`
: Repair the database. Renames the old database to suneido.db.bak

`-s[erver]`
: Run Suneido as a server.

`-v[ersion]`
: Display information about this version of Suneido. Similar to the 
[Built](<../Language/Reference/Built.md>) function.

`-w[eb][=port]`
: Run the web status monitor. This is automatic when running as a server. This option is for running standalone or as a client. If running as a client on the same computer as the server or standalone with multiple copies, then you need to specify a different port (on the client/standalone).