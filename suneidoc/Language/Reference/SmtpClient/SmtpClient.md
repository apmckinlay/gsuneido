#### SmtpClient

``` suneido
(smtp_server, user = false, password = false, 
    timeout = 60, helo_arg = "")
```

Starts an SMTP session with the given server.

The timeout argument is in seconds.  This is the amount of time the [SocketClient](<../SocketClient/SocketClient.md>) will wait for a response from the server before timing out. If the timeout is exceeded, an exception will be thrown. This timeout does not apply to the initial establishment of the socket connection.

For example:

``` suneido
sc = SmtpClient("mail.com");
sc.SendMail("me@mail.com", "me@mail.com", "mckinlay@suneido.com",
    "re. hello",
    "this is my message");
sc.Close();
```

If the server requires authentication, supply a user name and password.

**Note:** The current version of SmtpClient only handles AUTH LOGIN authentification.

Some servers also require a specific argument to the HELO/EHLO command.

See also: [SmtpSendMessage](<../SmtpSendMessage.md>), [PopClient](<../PopClient/PopClient.md>),
[RFC821 Simple Mail Transfer Protocol](<http://www.faqs.org/rfcs/rfc821.html>),
[RFC2554 SMTP Service Extension for Authentication](<http://www.faqs.org/rfcs/rfc2554.html>)