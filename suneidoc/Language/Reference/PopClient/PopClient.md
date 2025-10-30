#### PopClient

``` suneido
(pop_server, user, pass, timeout = 60)
```

Starts a POP session with the given server and authenticates using the  user and password given.  If there is an error during this process the connection is closed.  The timeout argument should be in seconds and defaults to 60 seconds.  This is the amount of time the [SocketClient](<../SocketClient/SocketClient.md>) will wait for a response from the server before timing out.  If the timeout is exceeded, an exception will be thrown. This timeout does not apply to the initial establishment of the socket connection.

For example:

``` suneido
pc = PopClient("mail.com", "andrew", "password")
for (i = 1; false isnt msg = pc.GetMessage(i); ++i)
    Print(msg)
pc.Close()
```

See also: 
[SmtpClient](<../SmtpClient/SmtpClient.md>)