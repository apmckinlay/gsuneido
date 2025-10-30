## Sending and Receiving Email

**Category:** Coding

**Problem**

You need to send or receive email messages from Suneido.

**Ingredients**

SmtpSendMessage, SmtpClient, and PopClient

**Recipe**

The standard method to send email over the internet is to use SMTP - Simple Mail Transport Protocol. The standard method to receive email is POP3 (Post Office Protocol version 3). There is also a more advanced system, IMAP (Internet Message Access Protocol), but there is no support for this yet.

*Note: The following examples use fictitious email addresses and servers. To actually run the examples, you'll need to substitute valid addresses and servers.*

<u>Sending messages using SmtpClient</u>

The easiest way to send an email using Suneido code is by using the SmtpSendMessage function, a simple wrapper for SmtpClient described below. SmtpSendMessage has the following format:

``` suneido
SmtpSendMessage(server, from, to, subject, message, header_extra = "")
```

The header_extra parameter is used for MIME headers and content other than plain text. The following examples don't use it.

For example, to send a message (assuming your server is mail.test.com and your account is jack@test.com and you want to send the message to john@smith.com), you could do the following:

``` suneido
SmtpSendMessage("mail.test.com", "jack@test.com", "john@smith.com",
    "test subject", "test message")
```

Note: In the past, many SMTP servers would accept email from anyone. However, these days, with spam issues, many SMTP servers will only accept mail from valid addresses for that mail server. Some SMTP servers also require authentication (i.e. a password), this isn't supported yet.

The message isn't limited to one line; it can be as long as you want (within reason). For example:

``` suneido
SmtpSendMessage("mail.test.com", "jack@test.com", "john@smith.com", "test subject",
"line one
line two
line three")
```

Of course, the message could also come from elsewhere, e.g. read from a file.

If you need to send more than one message at a time, or need more control over the header from and to, you can use SmtpClient. To start an SmtpClient:

``` suneido
smtp_client = SmptClient(server)
```

For example, to connect to the mail.test.com server, use:

``` suneido
smtp_client = SmtpClient("mail.test.com")
```

Then you can send any number of messages using the SendMail method:

``` suneido
smtp_client.SendMail(from, header_from, to, subject, message, header_extra = "", header_to = false)
```

'header_from' is used to specify what gets used for the From field of the email header, whereas 'from' is the actual email address the message is coming from. Normally, these will be the same. The same applies to 'to' and 'header_to'.  The 'header_extra' parameter is used for specifying different content types.  If you are only sending plain text then you don't need to use the 'header_extra' parameter.

For example:

``` suneido
smpt_client = SmtpClient("mail.test.com")
smtp_client.SendMail("jack@test.com", "jack@test.com", "jill@another.com", "Requested Docs", "The documents have been sent") 
// can send more emails here
smtp_client.Close()
```

The Close method of SmtpClient is used to close the connection.

<u>Receiving Mail Using PopClient</u>

PopClient can be used to retrieve messages, get a list of messages in the mailbox, get the size of messages in the mailbox, get a portion of a message from a mailbox and delete messages.

To connect to a POP3 server, create an instance of the PopClient class:

``` suneido
PopClient(server, user, pass)
```

For example:

``` suneido
pop_client = PopClient("mail.test.com", "jack@test.com", "test_password")
```

Note: Some POP3 servers expect the full email address for the user, and some require only the portion up to the '@', so in the latter case, the above example would become:

``` suneido
pop_client = PopClient("mail.test.com", "jack", "test_password")
```

To get a list of messages in the mailbox, use the List method:

``` suneido
list = pop_client.List()
```

The list will have two columns. The first is the message index and the second is the message size (in octets).

To get information about a message at a specific index in the mailbox:

``` suneido
list = pop_client.List(3) // gets info for 3rd message
```

To get the size of a message in the mailbox (size is in octets):

``` suneido
size = pop_client.GetMessageSize(3)  // gets size of 3rd message
```

To get a portion of a message in the mailbox, use the Top method:

``` suneido
msg_start = pop_client.Top(3, 15) // gets first 15 lines from 3rd message
```

To get an entire message:

``` suneido
msg = pop_client.GetMessage(3) // gets 3rd message,
```

GetMessage returns the message as a string, or false if there was an error getting the message.

To delete a message at a specific index in the mailbox:

``` suneido
result = pop_client.DeleteMessage(3) // deletes the 3rd message
```

The result will be true if successful, false otherwise.

To close the connection to the POP server, use the Close method:

``` suneido
pop_client.Close()
```

**See Also**

The specifications for these and other internet standards can be found in RFC's. [RFC 821](<http://www.faqs.org/rfcs/rfc821.html>) covers SMTP (sending email) and [RFC 1939](<http://www.faqs.org/rfcs/rfc1939.html>) covers POP3 (receiving email).

Suneido also includes simple implementations of an SMTP server (SmtpServer) and a POP3 server (Pop3Server). These are abstract base classes, you have to derive from them and implement the desired message storage.