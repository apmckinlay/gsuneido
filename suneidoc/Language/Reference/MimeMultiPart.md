<div style="float:right"><span class="toplinks"><a href="MimeMultiPart/Methods">Methods</a></span></div>

### MimeMultiPart

``` suneido
(subtype = "mixed")
```

A sub-class of [MimeBase](<MimeBase.md>) used for multipart or alternate formats.

For example:

``` suneido
MimeMultiPart('alternative').
    To("joe@hotmail.com").
    From('sue@mail.com').
    Subject('multipart test').
    Date().
    Message_ID().
    Attach(MimeText("hello world")).
    Attach(MimeText("<i>hello</i> <u>world</u>", 'html')).
    ToString()
```

would produce:

``` suneido
Content-Type: multipart/alternative; boundary="====================979610331"
MIME-Version: 1.0
To: joe@hotmail.com
From: sue@mail.com
Date: Fri, 19 Oct 2007 16:05:01 -0600
Subject: multipart test
Message-ID: <8a1bbfa4.d26c.4ce3.9267.698b5b0417ae@suneido.com>

--====================979610331
Content-Type: text/plain; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit

hello world
--====================979610331
Content-Type: text/html; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit

<b>hello</b> <big>world</big>
--====================979610331--
```