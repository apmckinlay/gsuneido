#### smtpclient.SendMail

``` suneido
(from, header_from, to, subject, message) => True or string
```

Sends the message to the recipient specified in the "to" parameter.  
"from" is the actual address of the sender.
"header_from" can be an alias and shows up on the header of the message.

Returns True if the send succeeds,
or else the string returned by the SMTP server.